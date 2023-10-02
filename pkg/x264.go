package pkg

import (
	"bytes"
	"errors"
	"image"
	"math"

	"github.com/gen2brain/x264-go"
)

type X264Encoder struct {
	Encoder *x264.Encoder
	buffer  *bytes.Buffer
}

func NewEncoder(size image.Rectangle) *X264Encoder {
	buffer := bytes.NewBuffer(make([]byte, 0))
	encoder, err := x264.NewEncoder(buffer, &x264.Options{
		Width:     int(math.Ceil(float64(size.Dx())/2)) * 2,
		Height:    int(math.Ceil(float64(size.Dy())/2)) * 2,
		FrameRate: 1,
		Tune:      "zerolatency",
		Preset:    "ultrafast",
		Profile:   "baseline",
	})

	if err != nil {
		return nil
	}

	return &X264Encoder{
		Encoder: encoder,
		buffer:  buffer,
	}
}

func (x *X264Encoder) Encode(frame *image.RGBA) ([]byte, error) {
	if x.Encoder == nil {
		return nil, errors.New("x264 encoder is nil!")
	}

	err := x.Encoder.Encode(frame)
	if err != nil {
		return nil, err
	}

	err = x.Encoder.Flush()
	if err != nil {
		return nil, err
	}

	frameBytes := x.buffer.Bytes()
	x.buffer.Reset()
	return frameBytes, nil
}
