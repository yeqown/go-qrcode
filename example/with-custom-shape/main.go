package main

import (
	qrcode "github.com/yeqown/go-qrcode"
)

type smallerCircle struct {
	smallerPercent float64
}

func (sc *smallerCircle) DrawFinder(ctx *qrcode.DrawContext) {
	backup := sc.smallerPercent
	sc.smallerPercent = 1.0
	sc.Draw(ctx)
	sc.smallerPercent = backup
}

func newShape(radiusPercent float64) qrcode.IShape {
	return &smallerCircle{smallerPercent: radiusPercent}
}

func (sc *smallerCircle) Draw(ctx *qrcode.DrawContext) {
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
	qrc, err := qrcode.New("with-custom-shape", qrcode.WithCustomShape(shape))
	// qrc, err := qrcode.New("with-custom-shape", qrcode.WithCircleShape())
	if err != nil {
		panic(err)
	}

	err = qrc.Save("./smaller.png")
	if err != nil {
		panic(err)
	}
}
