package main

import (
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"github.com/yeqown/go-qrcode/writer/standard/shapes"
)

func HChainBlock(ctx *standard.DrawContext) {
	w, h := ctx.Edge()
	fw, fh := float64(w), float64(h)
	x, y := ctx.UpperLeft()
	cx, cy := x+fw/2, y+fh/2
	r := fw * 0.85 / 2 // todo:
	l := r * 0.2

	ctx.SetColor(ctx.Color())

	mask := ctx.Neighbours()

	drawRect := func(x, y, w, h float64) {
		ctx.DrawRectangle(x, y, w, h)
		ctx.Fill()
	}
	_ = mask
	_ = drawRect

	ctx.DrawCircle(cx, cy, r)

	if mask&standard.NLeft|standard.NSelf == standard.NLeft|standard.NSelf {
		drawRect(x, cy-l, fw/2, 2*l)
	}
	if mask&standard.NRight|standard.NSelf == standard.NRight|standard.NSelf {
		drawRect(cx, cy-l, fw/2, 2*l)
	}

	ctx.Fill()
}

func main() {
	// assemble qr injecting build in function or you own for drawing

	shape := shapes.Assemble(shapes.RoundedFinder(), shapes.LiquidBlock())
	//shape := shapes.Assemble(shapes.RoundedFinder(), HChainBlock)

	qrc, err := qrcode.New(`https://github.com/yeqown/go-qrcode`)
	if err != nil {
		panic(err)
	}

	w, err := standard.New("./smaller.png",
		standard.WithCustomShape(shape),
	)
	if err != nil {
		panic(err)
	}

	err = qrc.Save(w)
	if err != nil {
		panic(err)
	}
}
