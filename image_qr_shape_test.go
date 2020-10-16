package qrcode

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"testing"
)

func Test_circle_Draw(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)
	rgba := image.NewRGBA(rect)
	// white background
	for posX := rect.Min.X; posX < rect.Max.X; posX++ {
		for posY := rect.Min.Y; posY < rect.Max.Y; posY++ {
			rgba.Set(posX, posY, color.White)
		}
	}

	ctx := &DrawContext{
		upperLeft:  rect.Min,
		lowerRight: rect.Max,
		img:        rgba,
		color:      color.White,
	}
	_shapeCircle.Draw(ctx)

	// save to file
	fd, err := os.Create("./testdata/circle.jpeg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = jpeg.Encode(fd, rgba, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
