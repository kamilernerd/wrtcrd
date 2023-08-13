package pkg

import (
	"fmt"
	"log"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

type Stream struct {
	Close chan int
	Peer  *webrtc.PeerConnection
}

type ChannelDeviceInput struct {
	Input  string `json:"input"`
	Button string `json:"button"`
	Action string `json:"action"`
}

func (s *Stream) NewWebrtcSession(sdp string, capturer *Capturer) (*webrtc.SessionDescription, error) {
	s.Peer = s.createPeer()
	s.Close = make(chan int, 1)

	for _, v := range capturer.Screens {
		track := s.addVideoTrack(v.Index)
		go s.writeTrack(track, v)
	}

	offer := webrtc.SessionDescription{
		SDP:  sdp,
		Type: webrtc.SDPTypeOffer,
	}

	// Set the remote SessionDescription
	if err := s.Peer.SetRemoteDescription(offer); err != nil {
		return nil, err
	}

	answer, err := s.Peer.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	if err = s.Peer.SetLocalDescription(answer); err != nil {
		return nil, err
	}

	return &answer, nil
}

func (s *Stream) createPeer() *webrtc.PeerConnection {
	m := webrtc.MediaEngine{}

	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264, ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        102,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))

	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun2.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		if connectionState == webrtc.ICEConnectionStateClosed {
			fmt.Println("ICE has closed the connection")
			s.Close <- 1
			return
		}

		if connectionState == webrtc.ICEConnectionStateDisconnected {
			fmt.Println("ICE has disconnected")
			s.Close <- 1
			return
		}

		if connectionState == webrtc.ICEConnectionStateFailed {
			fmt.Println("ICE Connection has gone to failed exiting")
			s.Close <- 1
			return
		}
	})

	peerConnection.OnConnectionStateChange(func(conn webrtc.PeerConnectionState) {
		if conn == webrtc.PeerConnectionStateConnected {
			return
		}

		// Closed due to some error, or just after a disconnect.
		if conn == webrtc.PeerConnectionStateFailed {
			fmt.Println("Peer Connection has gone to failed exiting")
			s.Close <- 1
			return
		}

		// Closed unexpectedly.
		if conn == webrtc.PeerConnectionStateClosed {
			fmt.Println("Peer has closed the connection")
			s.Close <- 1
			return
		}

		// Normal disconnect
		if conn == webrtc.PeerConnectionStateDisconnected {
			fmt.Println("Peer has disconnected")
			s.Close <- 1
			return
		}
	})

	return peerConnection
}

func (s *Stream) addVideoTrack(index int) *webrtc.TrackLocalStaticSample {
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
		MimeType: webrtc.MimeTypeH264,
	}, fmt.Sprintf("display-%d", index), fmt.Sprintf("remote-display-%d", index))

	rtpSender, err := s.Peer.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are retuned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()
	return videoTrack
}

func (s *Stream) writeTrack(track *webrtc.TrackLocalStaticSample, screen Screen) {
	encoder := NewEncoder(screen.Boundaries)

	for {
		select {
		case <-s.Close:
			return
		default:
			b, err := encoder.Encode(<-screen.Frame)
			if err != nil {
				log.Println("Failed to encode img in writeTrack")
				return
			}

			err = track.WriteSample(media.Sample{
				Data:     b,
				Duration: time.Duration(time.Millisecond),
			})

			if err != nil {
				log.Println("Frame lost, sorry :(")
				return
			}
		}
	}
}
