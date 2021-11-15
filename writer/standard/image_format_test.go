package standard

import (
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/fogleman/gg"
)

func newImage() image.Image {
	rect := image.Rect(0, 0, 100, 100)
	rgba := image.NewRGBA(rect)
	dc := gg.NewContextForRGBA(rgba)

	dc.DrawRectangle(0, 0, 100, 100)
	dc.SetColor(color.White)
	dc.Fill()

	dc.SetColor(color.Black)
	dc.DrawString("yeqown", 25, 50)

	return dc.Image()
}

func Test_JPEG_Encoder(t *testing.T) {
	img := newImage()

	fd, _ := os.OpenFile("./testdata/encoder_JPEG.jpeg", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	err := jpegEncoder{}.Encode(fd, img)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func Test_PNG_Encoder(t *testing.T) {
	img := newImage()

	fd, _ := os.OpenFile("./testdata/encoder_PNG.png", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	err := pngEncoder{}.Encode(fd, img)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
