package pkg

import (
	"image"
	"time"

	"github.com/kbinani/screenshot"
)

type Screen struct {
	Index      int
	Boundaries image.Rectangle
	Frame      chan *image.RGBA
}

type Capturer struct {
	Screens []Screen
}

func NewCapturer() *Capturer {
	return &Capturer{
		Screens: initCapturer(),
	}
}

func initCapturer() []Screen {
	monitorCount := screenshot.NumActiveDisplays()
	monitors := make([]Screen, monitorCount)

	for i := 0; i < monitorCount; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		monitors[i] = Screen{
			Index:      i,
			Frame:      make(chan *image.RGBA),
			Boundaries: bounds,
		}
	}
	return monitors
}

func (c *Capturer) FrameCapturer(index int) {
	delta := time.Duration(time.Second.Seconds()/60.0) * time.Millisecond

	for {
		screen := c.Screens[index]

		startedAt := time.Now()

		img, err := screenshot.CaptureRect(screen.Boundaries)
		if err != nil {
			return
		}

		screen.Frame <- img

		elapsed := time.Since(startedAt)
		sleepDuration := delta - elapsed

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}
