package beautifier

import "github.com/yeqown/go-qrcode/v2"

// Style defines how to render shapes.
type Style interface {
	// Draw renders a shape.
	// Returns the SVG content and true if the shape was handled.
	// Returns empty string and false if the shape was NOT handled (fallback needed).
	// ctx provides context like the full Matrix if needed for advanced rendering.
	Draw(ctx *DrawContext, shape Shape) (string, bool)
}

// DrawContext holds global context for drawing.
type DrawContext struct {
	Matrix qrcode.Matrix
}
