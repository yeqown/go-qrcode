package beautifier

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/yeqown/go-qrcode/v2"
)

// Writer implements qrcode.Writer interface for SVG output using Style-Shape architecture.
type Writer struct {
	style  Style
	closer io.WriteCloser
}

// New creates a new BeautifyWriter with the given Style.
func New(filename string, style Style) (*Writer, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, errors.Wrap(err, "create file failed")
	}
	return NewWithWriter(f, style), nil
}

// NewWithWriter creates a new BeautifyWriter with io.WriteCloser and Style.
func NewWithWriter(w io.WriteCloser, style Style) *Writer {
	return &Writer{
		style:  style,
		closer: w,
	}
}

func (w *Writer) Write(mat qrcode.Matrix) error {
	// 1. Analyze the matrix to find shapes
	shapes := Analyze(mat)
	
	// 2. Draw shapes using the Style
	ctx := &DrawContext{Matrix: mat}
	var content string
	
	for _, shape := range shapes {
		svg, handled := w.style.Draw(ctx, shape)
		if handled {
			content += svg
			continue
		}
		
		// Fallback: Decompose shape into smaller pieces
		// For now, we just decompose everything into single blocks.
		// A better approach would be to decompose Line -> Blocks, Square -> Blocks.
		// Since our base unit is Block, we just iterate points.
		for _, p := range shape.Points {
			blockShape := Shape{Type: ShapeTypeBlock, Points: []Point{p}}
			svg, handled := w.style.Draw(ctx, blockShape)
			if handled {
				content += svg
			} else {
				// Ultimate fallback: Draw a simple rect
				content += fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="1" fill="black" />`, p.X, p.Y)
			}
		}
	}

	// Construct the full SVG
	width := mat.Width()
	height := mat.Height()
	padding := 4
	fullWidth := width + padding*2
	fullHeight := height + padding*2

	svg := fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" version="1.1" viewBox="0 0 %d %d">
<rect width="100%%" height="100%%" fill="#ffffff"/>
<g transform="translate(%d, %d)">
%s
</g>
</svg>`,
		fullWidth, fullHeight,
		padding, padding,
		content,
	)

	_, err := w.closer.Write([]byte(svg))
	return err
}

func (w *Writer) Close() error {
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}
