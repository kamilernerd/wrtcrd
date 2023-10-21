package pkg

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

type MouseHandler struct{}

const (
	MouseMoveEvent   = "move"
	MouseActionEvent = "action"
)

func (m *MouseHandler) setMousePosition(dx int, dy int, x float64, y float64) {
	robotgo.Move(int(float64(dx)*x), int(float64(dy)*y))
}

func (m *MouseHandler) parseMouseInput(input string) (int64, float64, float64) {
	decodedMousePos, _ := base64.StdEncoding.DecodeString(input)

	mouseData := strings.Split(string(decodedMousePos), ";")

	d, _ := strconv.ParseInt(mouseData[0], 0, 64) // Display id
	x, _ := strconv.ParseFloat(mouseData[1], 64)  // x pos inside display preview
	y, _ := strconv.ParseFloat(mouseData[2], 64)  // y pos inside display preview

	return d, x, y
}

func (m *MouseHandler) MouseinputHandler(data []byte, c *Capturer) {
	var d Datagram
	json.Unmarshal(data, &d)

	switch d.Event {
	case MouseMoveEvent:
		d, x, y := m.parseMouseInput(d.Value.(string))
		dx := c.Screens[d].Boundaries.Dx()
		dy := c.Screens[d].Boundaries.Dy()

		// Corrects the mouse position on monitors coordinates
		// X pos screen 1 = 0.0
		// X pos screen 2 = 1.0...
		// The Y coordinate should always stay the same I guess.
		multiplier := 0.0
		if d >= int64(multiplier) {
			multiplier = float64(d)
		}

		m.setMousePosition(dx, dy, multiplier+x, y)
	case MouseActionEvent:
		if d.Value == 0 {
			robotgo.Toggle("left")
		} else if d.Value == 1 {
			robotgo.Toggle("center")
		} else if d.Value == 2 {
			robotgo.Toggle("right")
		} else { // Wheel scroll
			p := strings.Split(d.Value.(string), ";")
			x, _ := strconv.ParseInt(p[0], 0, 64)
			y, _ := strconv.ParseInt(p[1], 0, 64)
			robotgo.Scroll(int(x), int(y))
		}
	}
}
