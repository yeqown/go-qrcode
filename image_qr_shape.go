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
	// Draw to fill the IShape of qrcode.
	Draw(ctx *DrawContext)

	// DrawFinder to fill the finder pattern of QRCode, what's finder? google it for more information.
	DrawFinder(ctx *DrawContext)
}

// DrawContext is a rectangle area
type DrawContext struct {
	*gg.Context

	upperLeft image.Point // (x1, y1)
	w, h      int

	color color.Color
}

// UpperLeft returns the point which indicates the upper left position.
func (dc *DrawContext) UpperLeft() image.Point {
	return dc.upperLeft
}

// Edge returns width and height of each shape could take at most.
func (dc *DrawContext) Edge() (width, height int) {
	return dc.w, dc.h
}

// Color returns the color which should be fill into the shape. Note that if you're not
// using this color but your coded color.Color, some ImageOption functions those set foreground color
// would take no effect.
func (dc *DrawContext) Color() color.Color {
	return dc.color
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

func (r rectangle) DrawFinder(ctx *DrawContext) {
	r.Draw(ctx)
}

// circle IShape
type circle struct{}

// Draw
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

func (r circle) DrawFinder(ctx *DrawContext) {
	r.Draw(ctx)
}
