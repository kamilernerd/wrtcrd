import { useState } from "react";
import { KeyDownListener, KeyUpListener, MouseEventListener, MouseUpListener, MouseDownListener, MouseScrollListener } from "./lib";

type StreamsMap = Array<{
  stream: MediaStream,
  ref: HTMLVideoElement
  canvas: HTMLCanvasElement | null
}>

function App() {
  const FPS_DRAW_CAP = 1000 / 60 // 60fps cap
  const [streams, setStreams] = useState<StreamsMap>([])
  const socket = new WebSocket("ws://localhost:1337/ws");
  const pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: "stun:stun2.l.google.com:19302",
      }
    ],
  });

  pc.ondatachannel = (ev: RTCDataChannelEvent) => { }

  const keyboardDataChannel = pc.createDataChannel("keyboard", {
    negotiated: true,
    ordered: true,
    id: 1
  })

  const mouseDataChannel = pc.createDataChannel("mouse", {
    negotiated: true,
    ordered: true,
    id: 2
  })

  pc.onconnectionstatechange = async (ev: Event) => {
    console.log(pc.connectionState)
    if (pc.connectionState === "disconnected" || pc.connectionState === "failed" || pc.connectionState === "closed") {
      socket.close();

      // remove old streams
      const streams = document.getElementById("streams");
      streams!.innerHTML = '';
    }

    if (pc.connectionState === "connected") {
      console.info("Stats:", await pc.getStats())
      console.info("Configuration:", pc.getConfiguration())
    }
  }

  pc.ontrack = (evt: RTCTrackEvent) => {
    const videoElement = document.createElement("video")

    if (evt.streams[0].id === "audio-system") {
      videoElement.srcObject = evt.streams[0];
      videoElement.muted = false;
      videoElement.controls = true;
      videoElement.id = evt.streams[0].id
      videoElement.play();

      setStreams((prev) => [...prev, {
        stream: evt.streams[0],
        ref: videoElement,
        canvas: null
      }])

      document.querySelector("#streams")?.appendChild(videoElement)
    } else {
      const canvas = document.createElement("canvas")
      const ctx = canvas.getContext("2d")

      setStreams((prev) => [...prev, {
        stream: evt.streams[0],
        ref: videoElement,
        canvas: canvas,
      }])

      videoElement.srcObject = evt.streams[0];
      videoElement.muted = false;
      videoElement.controls = false;
      videoElement.style.display = "none"
      canvas.id = evt.streams[0].id.replace("remote-display-", "")

      videoElement.play();

      canvas.width = window.innerWidth / 2
      canvas.height = window.innerHeight / 2

      // Drawing loop
      videoElement.onplay = () => {
        function loop() {
          setTimeout(() => { }, FPS_DRAW_CAP) // SLEEP

          ctx?.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
          requestAnimationFrame(loop);
        }
        requestAnimationFrame(loop);
      }

      canvas.onkeydown = ev => KeyDownListener(ev, keyboardDataChannel)
      canvas.onkeyup = ev => KeyUpListener(ev, keyboardDataChannel)
      canvas.onmousemove = ev => MouseEventListener(ev, mouseDataChannel)
      canvas.onmousedown = ev => MouseDownListener(ev, mouseDataChannel)
      canvas.onmouseup = ev => MouseUpListener(ev, mouseDataChannel)
      canvas.onwheel = ev => MouseScrollListener(ev, mouseDataChannel)

      document.querySelector("#streams")?.appendChild(videoElement)
      document.querySelector("#streams")?.appendChild(canvas)
    }
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
    setTimeout(heartbeat, 100); // 100ms
  }

  socket.onmessage = function(e) {
    const message = JSON.parse(String(e.data)) as SocketMessage;
    switch (message.Event) {
      case "answer":
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
