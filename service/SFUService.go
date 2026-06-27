package service

import (
	"encoding/json"
	"log"
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
	log.Printf("[ws] new connection joined: %s", unsafeConn.RemoteAddr())
	c := &threadSafeWriter{unsafeConn, sync.Mutex{}}
	defer c.Close()

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Printf("[ws] failed to create peer connection for %s: %v", unsafeConn.RemoteAddr(), err)
		return
	}
	defer peerConnection.Close()

	if _, err := peerConnection.AddTransceiverFromKind(
		webrtc.RTPCodecTypeAudio,
		webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
	); err != nil {
		log.Printf("[ws] failed to add audio transceiver for %s: %v", unsafeConn.RemoteAddr(), err)
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
		log.Printf("[ws] peer connection state changed for %s: %s", unsafeConn.RemoteAddr(), p.String())

		switch p {
			case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Printf("Failed to close PeerConnection: %v", err)
			}
			case webrtc.PeerConnectionStateClosed:
				s.signalPeerConnections()
			default:
		}
	})

	peerConnection.OnTrack(func (t *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		log.Printf("[ws] new track received for %s: %s", unsafeConn.RemoteAddr(), t.ID())
		tracklocal := s.addTrack(t)
		if tracklocal == nil {
			log.Printf("[ws] failed to add track for %s", unsafeConn.RemoteAddr())
			return
		}
		defer s.removeTrack(tracklocal)

		rtpPkt := &rtp.Packet{}
		buf := make([]byte, 1500)
		for {
			log.Printf("tracks: %d", len(s.trackLocals))
			i, _, err := t.Read(buf)
			if err != nil {
				log.Printf("[ws] read failed for %s: %v", unsafeConn.RemoteAddr(), err)
				return
			}

			if err = rtpPkt.Unmarshal(buf[:i]); err != nil {
				log.Printf("[ws] failed to unmarshal RTP packet for %s: %v", unsafeConn.RemoteAddr(), err)
				return
			}
			rtpPkt.Extension = false
			rtpPkt.Extensions = nil;
			if err = tracklocal.WriteRTP(rtpPkt); err != nil {
				log.Printf("[ws] failed to write RTP packet for %s: %v", unsafeConn.RemoteAddr(), err)
				return
			}
		}
	})

	s.signalPeerConnections()

	message := &webSocketMessage{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			log.Printf("[ws] read failed for %s: %v", unsafeConn.RemoteAddr(), err)
			return
		}

		log.Printf("[ws] received message from %s: %s", unsafeConn.RemoteAddr(), string(raw))
		if err := json.Unmarshal(raw, message); err != nil {
			log.Printf("[ws] failed to unmarshal message from %s: %v", unsafeConn.RemoteAddr(), err)
			return
		}

		switch message.Event {
			case "candidate":
				candidate := webrtc.ICECandidateInit{}
				if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
					log.Printf("[ws] failed to unmarshal candidate from %s: %v", unsafeConn.RemoteAddr(), err)
					return
				}

				log.Printf("[ws] adding ICE candidate for %s: %v", unsafeConn.RemoteAddr(), candidate)
				if err := peerConnection.AddICECandidate(candidate); err != nil {
					log.Printf("[ws] failed to add ICE candidate for %s: %v", unsafeConn.RemoteAddr(), err)
					return
				}

			case "answer":
				answer := webrtc.SessionDescription{}
				if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
					log.Printf("[ws] failed to unmarshal answer from %s: %v", unsafeConn.RemoteAddr(), err)
					return
				}

				log.Printf("[ws] setting remote description for %s: %v", unsafeConn.RemoteAddr(), answer)
				if err := peerConnection.SetRemoteDescription(answer); err != nil {
					log.Printf("[ws] failed to set remote description for %s: %v", unsafeConn.RemoteAddr(), err)
					return
				}
			default:
				log.Printf("[ws] unknown event from message: %+v", message)
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
		log.Printf("[ws] failed to create local track for %s: %v", t.ID(), err)
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
						log.Printf("[ws] failed to remove track for %s: %v", s.peerConnections[i].websocketConn.RemoteAddr(), err)
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
						log.Printf("[ws] failed to add track for %s: %v", s.peerConnections[i].websocketConn.RemoteAddr(), err)
						return true
					}
				}
			}

			offer, err := s.peerConnections[i].peerConnection.CreateOffer(nil)
			if err != nil {
				log.Printf("[ws] failed to create offer for %s: %v", s.peerConnections[i].websocketConn.RemoteAddr(), err)
				return true
			}

			if err := s.peerConnections[i].peerConnection.SetLocalDescription(offer); err != nil {
				log.Printf("[ws] failed to set local description for %s: %v", s.peerConnections[i].websocketConn.RemoteAddr(), err)
				return true
			}

			offerString, _ := json.Marshal(offer)
			if err := s.peerConnections[i].websocketConn.WriteJSON(&webSocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				log.Printf("[ws] failed to send offer for %s: %v", s.peerConnections[i].websocketConn.RemoteAddr(), err)
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
