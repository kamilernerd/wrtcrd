package pkg

import (
	// "encoding/base64"
	"encoding/json"
	"fmt"
	// "strconv"
	// "strings"

	"github.com/pion/webrtc/v3"

	"github.com/go-vgo/robotgo"
)

type ChannelDeviceInput struct {
	Input  string `json:"input"`
	Button string `json:"button"`
	Action string `json:"action"`
}

func setMousePosition(dx int, dy int, x int, y int) {
	robotgo.Move(int(dx*x), int(dy*y))
}

func parseMouseInputEvent(input ChannelDeviceInput) {
	// decodedMousePos, _ := base64.StdEncoding.DecodeString(input.Button)
	//
	// mouseData := strings.Split(string(decodedMousePos), ";")
	//
	// x, _ := strconv.ParseFloat(mouseData[0], 64)
	// y, _ := strconv.ParseFloat(mouseData[1], 64)
	//
	// setMousePosition(dx, dy, x, y)
}

func MessageHandler(msg webrtc.DataChannelMessage, c *Capturer) {
	var d ChannelDeviceInput
	json.Unmarshal(msg.Data, d)

	// dx := c.Screens.Displays[displayNumber].Size.Dx()
	// dy := screen.Displays[displayNumber].Size.Dy()

	fmt.Println(d)
}
