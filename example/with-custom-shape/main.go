package main

import (
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type smallerCircle struct {
	smallerPercent float64
}

func (sc *smallerCircle) DrawFinder(ctx *standard.DrawContext) {
	backup := sc.smallerPercent
	sc.smallerPercent = 1.0
	sc.Draw(ctx)
	sc.smallerPercent = backup
}

func newShape(radiusPercent float64) standard.IShape {
	return &smallerCircle{smallerPercent: radiusPercent}
}

func (sc *smallerCircle) Draw(ctx *standard.DrawContext) {
	w, h := ctx.Edge()
	upperLeft := ctx.UpperLeft()
	color := ctx.Color()

	// choose a proper radius values
	radius := w / 2
	r2 := h / 2
	if r2 <= radius {
		radius = r2
	}

	// 80 percent smaller
	radius = int(float64(radius) * sc.smallerPercent)

	cx, cy := upperLeft.X+w/2, upperLeft.Y+h/2 // get center point
	ctx.DrawCircle(float64(cx), float64(cy), float64(radius))
	ctx.SetColor(color)
	ctx.Fill()

}

func main() {
	shape := newShape(0.7)
	qrc, err := qrcode.New("with-custom-shape")
	// qrc, err := qrcode.New("with-custom-shape", qrcode.WithCircleShape())
	if err != nil {
		panic(err)
	}

	w, err := standard.New("./smaller.png", standard.WithCustomShape(shape))
	if err != nil {
		panic(err)
	}

	err = qrc.Save(w)
	if err != nil {
		panic(err)
	}
}
