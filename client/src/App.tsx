import { useState } from "react";
import { KeyDownListener, KeyUpListener } from "./lib";

function App() {
  const [streams, setStreams] = useState<Array<{
    stream: MediaStream,
    ref: HTMLVideoElement
  }>>([])
  const socket = new WebSocket("ws://localhost:1337/ws");
  const pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: "stun:stun2.l.google.com:19302",
      }
    ],
  });

  pc.ondatachannel = (ev: RTCDataChannelEvent) => {
    console.log(ev)
  }

  const keyboardDataChannel = pc.createDataChannel("keyboard", {
    negotiated: true,
    ordered: true,
    id: 1
  })

  pc.onconnectionstatechange = async (ev: Event) => {
    console.log(pc.connectionState)
    if (pc.connectionState === "disconnected" || pc.connectionState === "failed" || pc.connectionState === "closed") {
      socket.close();
    }

    if (pc.connectionState === "connected") {
      console.info("Stats:", await pc.getStats())
      console.info("Configuration:", pc.getConfiguration())
    }
  }

  pc.ontrack = (evt: RTCTrackEvent) => {
    console.log(evt.streams);
    const videoElement = document.createElement("video")
    setStreams((prev) => [...prev, {
      stream: evt.streams[0],
      ref: videoElement,
    }])

    videoElement.width = 640
    videoElement.height = 480
    videoElement.srcObject = evt.streams[0];
    videoElement.muted = false;
    videoElement.controls = true;
    videoElement.play();

    videoElement.onkeydown = ev => KeyDownListener(ev, keyboardDataChannel)
    videoElement.onkeyup = ev => KeyUpListener(ev, keyboardDataChannel)

    document.querySelector("#streams")?.appendChild(videoElement)
  }

  socket.onopen = async (ev: Event) => {
    pc.restartIce()
    socket.send(JSON.stringify({
      Event: "offer",
      Value: await createOffer(),
    }))
    heartbeat();
  }

  socket.onclose = (e) => {
    keyboardDataChannel.close();
    pc.close();
  }

  function heartbeat() {
    socket.send("heartbeat");
    setTimeout(heartbeat, 200); // 200ms
  }

  socket.onmessage = function(e) {
    const message = JSON.parse(String(e.data)) as SocketMessage;
    switch (message.Event) {
      case "answer":
        console.log(message);
        pc.setRemoteDescription(new RTCSessionDescription(message.Value));
        break;
    }
  }

  const createOffer = async (): Promise<string | undefined> => {
    return new Promise((accept, reject) => {
      pc.onicecandidate = (evt) => {
        if (!evt.candidate) {
          const { sdp: offer } = pc.localDescription!;
          accept(offer);
        }
      };

      // Offer to receive multiple these tracks
      pc.addTransceiver("video")
      pc.addTransceiver("video")
      pc.addTransceiver("audio")

      pc.createOffer({
        iceRestart: true,
        offerToReceiveAudio: true,
        offerToReceiveVideo: true,
      })
        .then((ld) => {
          pc.setLocalDescription(ld);
        })
        .catch(reject);
    });
  };

  return (
    <>
      <div id="streams">
        <>
          {streams.length === 0 ?? (
            <p>No active streams</p>
          )}
        </>
      </div>
      <button onClick={async () => {
        socket.send(JSON.stringify({
          Event: "stop",
          Value: null,
        }));
      }}>STOP</button>
    </>
  )
}

export default App
