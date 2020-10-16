package qrcode

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
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
	*gg.Context

	upperLeft  image.Point // (x1, y1)
	lowerRight image.Point // (x2, y2)

	color color.Color
}

// rectangle IShape
type rectangle struct{}

func (r rectangle) Draw(c *DrawContext) {
	c.DrawRectangle(float64(c.upperLeft.X), float64(c.upperLeft.Y),
		float64(c.lowerRight.X), float64(c.lowerRight.Y))
	c.SetColor(c.color)
	c.Fill()
}

// circle IShape
type circle struct{}

// FIXED: Draw could not draw circle
func (r circle) Draw(c *DrawContext) {
	w := c.lowerRight.X - c.upperLeft.X
	h := c.lowerRight.Y - c.upperLeft.Y

	// choose property radius values
	radius := w / 2
	r2 := h / 2
	if r2 <= radius {
		radius = r2
	}

	cx, cy := c.upperLeft.X+w/2, c.upperLeft.Y+h/2 // get center point
	c.DrawCircle(float64(cx), float64(cy), float64(radius))
	c.SetColor(c.color)
	c.Fill()
}
