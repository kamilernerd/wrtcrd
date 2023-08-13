import { useRef } from "react";

function App() {
  const videoRef = useRef<HTMLVideoElement | null>();
  const mediaStream = new MediaStream()
  const socket = new WebSocket("ws://localhost:8080/ws");
  const pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: "stun:stun2.l.google.com:19302",
      }
    ],
  })

  pc.ontrack = (evt: RTCTrackEvent) => {
    console.info("ontrack triggered");
    if (videoRef.current) {
      evt.streams[0].getVideoTracks()
      // videoRef.current.srcObject = evt.streams[0];
      // videoRef.current.muted = true;
      // videoRef.current.play();
    }
  };

  socket.onopen = async (ev: Event) => {
    socket.send(JSON.stringify({
      Event: "offer",
      Value: await createOffer()
    }))

    heartbeat()
  }

  socket.onclose = (e) => {

  }

  function heartbeat() {
    socket.send("heartbeat");
    setTimeout(heartbeat, 15000); // 15 sec
  }

  socket.onmessage = function (e) {
    const message = JSON.parse(String(e.data)) as SocketMessage
    switch (message.Event) {
      case "answer":
        console.log(message)
        pc.setRemoteDescription(new RTCSessionDescription(message.Value));
        break
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

      pc.createOffer({
        offerToReceiveAudio: false,
        offerToReceiveVideo: true
      })
        .then((ld) => {
          pc.setLocalDescription(ld);
        })
        .catch(reject);
    });
  };

  return (
    <>
      <button onClick={async () => {
        socket.send(JSON.stringify({
          Event: "stop",
          Value: null
        }))
      }}>STOP</button>
      <video width="640" height="480" ref={(r) => (videoRef.current = r)} controls autoPlay muted />
    </>
  )
}

export default App
