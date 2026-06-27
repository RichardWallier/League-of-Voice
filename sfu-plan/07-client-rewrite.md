# Step 7 — Rewrite the client (audio-only SFU)

> [SFU plan](../SFU_PLAN.md) · Step 7 of 8 · prev: [06](./06-signal-renegotiation.md) · next: [08](./08-run-and-verify.md)

**Goal:** `test-client.html` speaks the SFU protocol: **server offers, client answers**; audio-only.

**Why:** The mesh client (peers offer each other, `from`/`to`, higher-ID-initiates) is
obsolete. The SFU client is simpler — it only answers and trickles candidates. Base it on
`sfu-ws/index.html` (~100 lines), trimmed to audio.

**Files:** `test-client.html`.

**Do:** replace the script body with this shape (keep your log box / buttons if you like):
```js
const stream = await navigator.mediaDevices.getUserMedia({ audio: true, video: false });
const pc = new RTCPeerConnection();
stream.getTracks().forEach(t => pc.addTrack(t, stream));

const ws = new WebSocket('ws://localhost:3000/ws');

pc.onicecandidate = e => {
  if (e.candidate) ws.send(JSON.stringify({ event: 'candidate', data: JSON.stringify(e.candidate) }));
};

pc.ontrack = e => {                       // audio-only: one <audio> per remote stream
  const el = document.createElement('audio');
  el.srcObject = e.streams[0];
  el.autoplay = true;
  document.body.appendChild(el);
};

ws.onmessage = async (evt) => {
  const msg = JSON.parse(evt.data);
  if (msg.event === 'offer') {
    await pc.setRemoteDescription(JSON.parse(msg.data));
    const answer = await pc.createAnswer();
    await pc.setLocalDescription(answer);
    ws.send(JSON.stringify({ event: 'answer', data: JSON.stringify(answer) }));
  } else if (msg.event === 'candidate') {
    await pc.addIceCandidate(JSON.parse(msg.data));
  }
};
```

Differences from the example's `index.html`: `getUserMedia` is audio-only, and `ontrack`
builds an `<audio>` element instead of `<video>`. **The client never calls `createOffer`** —
the server drives negotiation.

**Verify:** page loads, grabs the mic, connects without errors. The real end-to-end test is step 8.

**Next:** [Step 8 — run & verify](./08-run-and-verify.md).
