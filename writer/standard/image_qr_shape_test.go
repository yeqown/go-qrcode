package standard

import (
	"image"
	"image/color"
	"testing"

	"github.com/fogleman/gg"
)

func Test_rectangle_Draw(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)
	rgba := image.NewRGBA(rect)
	dc := gg.NewContextForRGBA(rgba)

	dc.DrawRectangle(0, 0, 100, 100)
	dc.SetColor(color.White)
	dc.Fill()

	ctx := &DrawContext{
		Context: dc,
		x:       0.0,
		y:       0.0,
		w:       50,
		h:       50,
		color:   color.Black,
	}
	_shapeRectangle.Draw(ctx)

	err := dc.SavePNG("./testdata/rectangle.png")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func Test_circle_Draw(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)
	rgba := image.NewRGBA(rect)
	dc := gg.NewContextForRGBA(rgba)

	dc.DrawRectangle(0, 0, 100, 100)
	dc.SetColor(color.White)
	dc.Fill()

	ctx := &DrawContext{
		Context: dc,
		x:       0.0,
		y:       0.0,
		w:       50,
		h:       50,
		color:   color.Black,
	}
	_shapeCircle.Draw(ctx)

	err := dc.SavePNG("./testdata/circle.png")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func Test_gg(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)
	rgba := image.NewRGBA(rect)
	dc := gg.NewContextForRGBA(rgba)

	dc.DrawRectangle(0, 0, 100, 100)
	dc.SetColor(color.White)
	dc.Fill()

	dc.DrawCircle(50, 50, 40)
	dc.SetColor(color.Black)
	dc.Fill()
	_ = dc.SavePNG("./testdata/out.png")
}
