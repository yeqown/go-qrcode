package qrcode

import (
	"image"
	"image/color"
	"math"
)

var (
	_shapeRectangle IShape = rectangle{}
	_shapeCircle    IShape = circle{}
)

type IShape interface {
	// draw to fill the IShape
	Draw(ctx *DrawContext)
}

// DrawContext is a rectangle area
type DrawContext struct {
	upperLeft  image.Point // (x1, y1)
	lowerRight image.Point // (x2, y2)

	img   *image.RGBA
	color color.Color
}

func (ctx *DrawContext) Set(x, y int) {
	if ctx == nil || ctx.img == nil {
		return
	}

	ctx.img.Set(x, y, ctx.color)
}

// rectangle IShape
type rectangle struct{}

func (r rectangle) Draw(ctx *DrawContext) {
	for posX := ctx.upperLeft.X; posX < ctx.lowerRight.X; posX++ {
		for posY := ctx.upperLeft.Y; posY < ctx.lowerRight.Y; posY++ {
			ctx.Set(posX, posY)
		}
	}
}

// circle IShape
type circle struct{}

// FIXME: Draw could not draw circle
func (r circle) Draw(ctx *DrawContext) {

	w := ctx.lowerRight.X - ctx.upperLeft.X
	h := ctx.lowerRight.Y - ctx.upperLeft.Y

	// choose property radius value
	radius := w / 2
	r2 := h / 2
	if r2 <= radius {
		radius = r2
	}

	cx, cy := ctx.lowerRight.X+w/2, ctx.lowerRight.Y+h/2 // get center point

	// Draw x,y
	for ; radius > 0; radius-- {
		r_2 := radius * radius
		for x := 0; x < radius; x++ {
			x_2 := x * x

			y := int(math.Sqrt(float64(r_2 - x_2)))
			ctx.Set(cx+x, cx-y)
			ctx.Set(cx+x, cy+y)
			ctx.Set(cx-x, cx-y)
			ctx.Set(cx-x, cy+y)
		}
	}
}
