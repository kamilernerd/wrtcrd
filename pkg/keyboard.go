package pkg

import (
	"encoding/json"
	"runtime"
	"strings"

	"github.com/go-vgo/robotgo"
)

type KeyboardHandler struct {
	keystate map[string]int
}

type KeyboardInput struct {
	State int
	Key   string
	Shift bool
	Meta  bool
	Ctrl  bool
	Alt   bool
}

const (
	KeyUP   = 0
	KeyDown = 1
)

const (
	windows = "windows"
	darwin  = "darwin"
	linux   = "linux"
)

type KeyboardDatagram struct {
	Datagram
	Value KeyboardInput
}

func (k *KeyboardHandler) UnmarshallDatagram(data []byte) KeyboardInput {
	var j KeyboardDatagram
	json.Unmarshal(data, &j)
	return j.Value
}

func (k *KeyboardHandler) isPlatform(platform string) bool {
	return runtime.GOOS == platform
}

func (k *KeyboardHandler) parseKeys(value KeyboardInput) KeyboardInput {
	switch value.Key {
	case "metaleft":
		if k.isPlatform(windows) {
			value.Key = "win"
		} else if k.isPlatform(darwin) {
			value.Key = "lcmd"
		} else {
			value.Key = "lmeta"
		}
	case "metaright":
		if k.isPlatform(windows) {
			value.Key = "win"
		} else if k.isPlatform(darwin) {
			value.Key = "rcmd"
		} else {
			value.Key = "rmeta"
		}
	case "altleft":
		value.Key = "lalt"
	case "altright":
		value.Key = "ralt"
	case "controlleft":
		value.Key = "lctrl"
	case "controlright":
		value.Key = "rctrl"
	case "shiftleft":
		value.Key = "lshift"
	case "shiftright":
		value.Key = "rshift"
	case "arrowup":
		value.Key = "up"
	case "arrowdown":
		value.Key = "down"
	case "arrowleft":
		value.Key = "left"
	case "arrowright":
		value.Key = "right"
	case "numlock":
		value.Key = "num_lock"
	case "numpadequal":
		value.Key = "num_equal"
	case "numpaddivide":
		value.Key = "num/"
	case "numpadmultiply":
		value.Key = "num*"
	case "numpadsubtrack":
		value.Key = "num-"
	case "numpadenter":
		value.Key = "num_enter"
	case "numpaddecimal":
		value.Key = "num."
	default:
		// handle normal keys a,b,c,d etc
		if strings.Contains(value.Key, "key") {
			value.Key = strings.Replace(value.Key, "key", "", len(value.Key))

			if value.Shift {
				value.Key = strings.ToUpper(value.Key)
			} else {
				value.Key = strings.ToLower(value.Key)
			}
		}

		// Handle numerical values
		if strings.Contains(value.Key, "numpad") {
			value.Key = strings.Replace(value.Key, "numpad", "", len(value.Key))
		} else if strings.Contains(value.Key, "digit") {
			value.Key = strings.Replace(value.Key, "digit", "", len(value.Key))
		}
	}
	return value
}

func (c *KeyboardHandler) HandleKeyboardInput(value KeyboardInput) {
	parsedKey := c.parseKeys(value)

	if parsedKey.State == KeyUP {
		robotgo.KeyToggle(parsedKey.Key, "up")
	} else if parsedKey.State == KeyDown {
		robotgo.KeyToggle(parsedKey.Key)
	}
}
