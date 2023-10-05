package pkg

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Capturer *Capturer
	Stream   *Stream
}

type Datagram struct {
	Event string
	Value interface{}
}

func NewClient() *Client {
	return &Client{
		Conn:     nil,
		Capturer: NewCapturer(),
		Stream:   &Stream{},
	}
}

func (c *Client) Close() {
	c.Conn.Close()
	c.Stream.Peer.Close()
	c.Conn = nil
}

func (c *Client) NewConnection(conn *websocket.Conn) {
	defer c.Close()

	c.Conn = conn

	for _, s := range c.Capturer.Screens {
		go c.Capturer.FrameCapturer(s.Index)
	}

	c.Send(Datagram{
		Event: "connect",
		Value: "connected",
	})

	for {
		d, err := c.Read()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseInternalServerErr) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch d.Event {
		case "offer":
			answer, err := c.Stream.NewWebrtcSession(d.Value.(string), c.Capturer)
			if err != nil {
				fmt.Println(err)
				break
			}
			c.Send(Datagram{
				Event: "answer",
				Value: answer,
			})
		case "stop":
			c.Close()
			break
		}
	}
}

func (c *Client) Read() (*Datagram, error) {
	d := &Datagram{}

	_, message, err := c.Conn.ReadMessage()
	if err != nil {
		log.Println("Error during message reading:", err)
		return nil, err
	}

	json.Unmarshal(message, d)
	return d, nil
}

func (c *Client) Send(d Datagram) {
	b, err := json.Marshal(d)
	if err != nil {
		log.Println("Error during message writing:", err)
		return
	}

	err = c.Conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("Error during message writing:", err)
		return
	}
}
