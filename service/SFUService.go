package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type peerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	websocketConn      *threadSafeWriter
}

type SFUService struct {
	listLock	sync.RWMutex
	peerConnections []peerConnectionState
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type webSocketMessage struct {
	Event string `json:"event"`
	Data string `json:"data"`
}

// Helper to make Gorilla Websockets threadsafe.
type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func NewSFUService() *SFUService {
	return &SFUService{
		trackLocals:     make(map[string]*webrtc.TrackLocalStaticRTP),
	}
}

func (t *threadSafeWriter) WriteJSON(v any) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

func (s *SFUService) Join(unsafeConn *websocket.Conn) {
	c := &threadSafeWriter{unsafeConn, sync.Mutex{}}
	defer c.Close()

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return
	}
	defer peerConnection.Close()

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()
	go func() {
	    for range pingTicker.C {
	        if err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
	            return // connection gone — stop pinging
	        }
	    }
	}()

	if _, err := peerConnection.AddTransceiverFromKind(
		webrtc.RTPCodecTypeAudio,
		webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
	); err != nil {
		return
	}

	s.listLock.Lock()
	s.peerConnections = append(s.peerConnections, peerConnectionState{peerConnection, c})
	s.listLock.Unlock()

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		b, _ := json.Marshal(i.ToJSON())
		c.WriteJSON(&webSocketMessage{
			Event: "candidate",
			Data:  string(b),
		})
	})

	peerConnection.OnConnectionStateChange(func (p webrtc.PeerConnectionState) {
		switch p {
			case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				return
			}
			case webrtc.PeerConnectionStateClosed:
				s.signalPeerConnections()
			default:
		}
	})

	peerConnection.OnTrack(func (t *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		tracklocal := s.addTrack(t)
		if tracklocal == nil {
			return
		}
		defer s.removeTrack(tracklocal)

		rtpPkt := &rtp.Packet{}
		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if err = rtpPkt.Unmarshal(buf[:i]); err != nil {
				return
			}
			rtpPkt.Extension = false
			rtpPkt.Extensions = nil;
			if err = tracklocal.WriteRTP(rtpPkt); err != nil {
				return
			}
		}
	})

	s.signalPeerConnections()

	message := &webSocketMessage{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			return
		}

		if err := json.Unmarshal(raw, message); err != nil {
			return
		}

		switch message.Event {
			case "candidate":
				candidate := webrtc.ICECandidateInit{}
				if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
					return
				}

				if err := peerConnection.AddICECandidate(candidate); err != nil {
					return
				}

			case "answer":
				answer := webrtc.SessionDescription{}
				if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
					return
				}

				if err := peerConnection.SetRemoteDescription(answer); err != nil {
					return
				}
			default:
		}
	}
	// end for loop
}

func (s *SFUService) addTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP  {
	s.listLock.Lock()
	defer func () {
		s.listLock.Unlock()
		s.signalPeerConnections()
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		return nil
	}
	s.trackLocals[t.ID()] = trackLocal

	return trackLocal
}

func (s *SFUService) removeTrack(t *webrtc.TrackLocalStaticRTP){
	s.listLock.Lock()
	defer func () {
		s.listLock.Unlock()
		s.signalPeerConnections()
	}()

	delete(s.trackLocals, t.ID())
}

func (s *SFUService) signalPeerConnections() {
	s.listLock.Lock()
	defer s.listLock.Unlock()

	attemptSync := func () (tryAgain bool) {
		for i := range s.peerConnections {
			if s.peerConnections[i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				 s.peerConnections = append(s.peerConnections[:i], s.peerConnections[i+1:]...)
				return true
			}
			existingSenders := map[string]bool{}
			for _, sender := range s.peerConnections[i].peerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true
				if _, ok := s.trackLocals[sender.Track().ID()]; !ok {
					if err := s.peerConnections[i].peerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			for _, receiver := range s.peerConnections[i].peerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			//check this later,  didnt understood functionality of this code, but it seems to be adding tracks to the peer connection if they are not already present
			for _, trackLocal := range s.trackLocals {
				if _, ok := existingSenders[trackLocal.ID()]; !ok {
					if _, err := s.peerConnections[i].peerConnection.AddTrack(trackLocal); err != nil {
						return true
					}
				}
			}

			offer, err := s.peerConnections[i].peerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err := s.peerConnections[i].peerConnection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, _ := json.Marshal(offer)
			if err := s.peerConnections[i].websocketConn.WriteJSON(&webSocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				return true
			}
		}
		return tryAgain
	}


	// not understood this loop
	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				s.signalPeerConnections()
			}()

			return
		}

		if !attemptSync() {
			break
		}
	}
}
