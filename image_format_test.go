package qrcode

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

	fd, _ := os.OpenFile("./testdata/JPEG_encoder_test.jpeg", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	err := jpegEncoder{}.Encode(fd, img)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func Test_PNG_Encoder(t *testing.T) {
	img := newImage()

	fd, _ := os.OpenFile("./testdata/PNG_encoder_test.png", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	err := pngEncoder{}.Encode(fd, img)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func Test_HEIF_Encoder(t *testing.T) {
	t.Skip("not implemented")

	img := newImage()

	fd, _ := os.OpenFile("./testdata/HEIF_encoder_test.heif", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	err := heifEncoder{}.Encode(fd, img)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
