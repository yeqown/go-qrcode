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

	upperLeft image.Point // (x1, y1)
	w, h      int

	color color.Color
}

// rectangle IShape
type rectangle struct{}

func (r rectangle) Draw(c *DrawContext) {
	// FIXED(@yeqown): miss parameter of DrawRectangle
	c.DrawRectangle(float64(c.upperLeft.X), float64(c.upperLeft.Y),
		float64(c.w), float64(c.h))
	c.SetColor(c.color)
	c.Fill()
}

// circle IShape
type circle struct{}

// FIXED: Draw could not draw circle
func (r circle) Draw(c *DrawContext) {
	// choose a proper radius values
	radius := c.w / 2
	r2 := c.h / 2
	if r2 <= radius {
		radius = r2
	}

	cx, cy := c.upperLeft.X+c.w/2, c.upperLeft.Y+c.h/2 // get center point
	c.DrawCircle(float64(cx), float64(cy), float64(radius))
	c.SetColor(c.color)
	c.Fill()
}
