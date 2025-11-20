package beautifier

import (
	"fmt"
	"strings"
)

// Neighbor direction constants
const (
	NTop    = 1 << 0
	NRight  = 1 << 1
	NBottom = 1 << 2
	NLeft   = 1 << 3
)

// getNeighbours returns a bitmask of which neighbors are set
func getNeighbours(bitmap [][]bool, x, y int) int {
	var neighbors int
	height := len(bitmap)
	if height == 0 {
		return 0
	}
	width := len(bitmap[0])

	// Top
	if y > 0 && bitmap[y-1][x] {
		neighbors |= NTop
	}
	// Right
	if x < width-1 && bitmap[y][x+1] {
		neighbors |= NRight
	}
	// Bottom
	if y < height-1 && bitmap[y+1][x] {
		neighbors |= NBottom
	}
	// Left
	if x > 0 && bitmap[y][x-1] {
		neighbors |= NLeft
	}

	return neighbors
}


// BaseStyle is a simple style that renders everything as standard blocks/shapes.
type BaseStyle struct {
	Color string
}

func (s BaseStyle) Draw(ctx *DrawContext, shape Shape) (string, bool) {
	switch shape.Type {
	case ShapeTypeFinder:
		// Draw standard finder
		// We can reuse the logic from previous FinderPattern
		// Assuming shape.Points[0] is top-left
		p := shape.Points[0]
		return fmt.Sprintf(`<g fill="%s"><path d="M%d %d h7 v7 h-7 z M%d %d v5 h5 v-5 z" /><rect x="%d" y="%d" width="3" height="3" /></g>`, 
			s.Color, p.X, p.Y, p.X+1, p.Y+1, p.X+2, p.Y+2), true
			
	case ShapeTypeBlock, ShapeTypeAlignment, ShapeTypeTiming:
		p := shape.Points[0]
		return fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="1" fill="%s" />`, p.X, p.Y, s.Color), true
		
	case ShapeTypeLineX:
		p := shape.Points[0]
		length := len(shape.Points)
		return fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="1" rx="0.5" ry="0.5" fill="%s" />`, p.X, p.Y, length, s.Color), true
		
	case ShapeTypeLineY:
		p := shape.Points[0]
		length := len(shape.Points)
		return fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="%d" rx="0.5" ry="0.5" fill="%s" />`, p.X, p.Y, length, s.Color), true
		
	case ShapeTypeSquare:
		// Draw circle for 2x2
		p := shape.Points[0]
		return fmt.Sprintf(`<circle cx="%d" cy="%d" r="1" fill="%s" />`, p.X+1, p.Y+1, s.Color), true
		
	case ShapeTypeL:
		p := shape.Points[0]
		return fmt.Sprintf(`<path d="M%d %d h2 v1 h-1 v1 h-1 z" fill="%s" />`, p.X, p.Y, s.Color), true
	}
	
	return "", false
}

// CompositeStyle allows combining multiple styles.
// It tries styles in order.
type CompositeStyle struct {
	Styles []Style
}

func (s CompositeStyle) Draw(ctx *DrawContext, shape Shape) (string, bool) {
	for _, style := range s.Styles {
		svg, handled := style.Draw(ctx, shape)
		if handled {
			return svg, true
		}
	}
	return "", false
}

// FuncStyle allows defining a style with a function for specific shapes.
type FuncStyle struct {
	Fn func(ctx *DrawContext, shape Shape) (string, bool)
}

func (s FuncStyle) Draw(ctx *DrawContext, shape Shape) (string, bool) {
	return s.Fn(ctx, shape)
}

// LiquidStyle creates fluid, connected blocks with rounded corners.
type LiquidStyle struct {
	Color string
}

func (s LiquidStyle) Draw(ctx *DrawContext, shape Shape) (string, bool) {
	// LiquidStyle handles all shape types by drawing each point with connections
	if len(shape.Points) == 0 {
		return "", false
	}
	
	// Special handling for Finder patterns - render as concentric circles
	if shape.Type == ShapeTypeFinder {
		return s.drawFinderCircle(shape.Points[0].X, shape.Points[0].Y), true
	}
	
	var sb strings.Builder
	bitmap := ctx.Matrix.Bitmap()
	
	// For each point in the shape, draw it with liquid connections
	for _, p := range shape.Points {
		x, y := p.X, p.Y
		
		// Get neighbors
		neighbors := getNeighbours(bitmap, x, y)
		
		// Draw the block with rounded corners based on neighbors
		sb.WriteString(s.drawLiquidBlock(x, y, neighbors))
	}
	
	return sb.String(), true
}

func (s LiquidStyle) drawFinderCircle(x, y int) string {
	// Draw Finder pattern as concentric circles (bullseye style)
	cx := float64(x) + 3.5
	cy := float64(y) + 3.5
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<g fill="%s">`, s.Color))
	
	// Outer ring (radius 3.5, hole radius 2.5)
	sb.WriteString(fmt.Sprintf(`<path d="M %.2f %.2f A 3.5 3.5 0 1 0 %.2f %.2f A 3.5 3.5 0 1 0 %.2f %.2f Z M %.2f %.2f A 2.5 2.5 0 1 1 %.2f %.2f A 2.5 2.5 0 1 1 %.2f %.2f Z" />`,
		cx, cy-3.5, cx, cy+3.5, cx, cy-3.5, // Outer circle
		cx, cy-2.5, cx, cy+2.5, cx, cy-2.5, // Inner hole (reverse winding)
	))
	
	// Inner dot (radius 1.5)
	sb.WriteString(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="1.5" />`, cx, cy))
	
	sb.WriteString("</g>")
	return sb.String()
}


func (s LiquidStyle) drawLiquidBlock(x, y int, neighbors int) string {
	// Match the standard LiquidBlock implementation
	fx := float64(x)
	fy := float64(y)
	cx := fx + 0.5  // Center X
	cy := fy + 0.5  // Center Y
	r := 0.5        // Radius (half of block size)
	
	var sb strings.Builder
	
	// Helper to check if mask has all specified bits
	has := func(mask, bits int) bool {
		return mask&bits == bits
	}
	
	// Draw horizontal/vertical connection bars
	if has(neighbors, NLeft|NRight) {
		// Horizontal bar through center
		sb.WriteString(fmt.Sprintf(`<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" fill="%s" />`,
			fx-0.5, cy-r, 2.0, 2*r, s.Color))
	}
	if has(neighbors, NTop|NBottom) {
		// Vertical bar through center
		sb.WriteString(fmt.Sprintf(`<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" fill="%s" />`,
			cx-r, fy-0.5, 2*r, 2.0, s.Color))
	}
	
	// Draw individual direction connectors
	if has(neighbors, NLeft) {
		sb.WriteString(fmt.Sprintf(`<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" fill="%s" />`,
			fx, cy-r, 0.5, 2*r, s.Color))
	}
	if has(neighbors, NRight) {
		sb.WriteString(fmt.Sprintf(`<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" fill="%s" />`,
			cx, cy-r, 0.5, 2*r, s.Color))
	}
	if has(neighbors, NTop) {
		sb.WriteString(fmt.Sprintf(`<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" fill="%s" />`,
			cx-r, fy, 2*r, 0.5, s.Color))
	}
	if has(neighbors, NBottom) {
		sb.WriteString(fmt.Sprintf(`<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" fill="%s" />`,
			cx-r, cy, 2*r, 0.5, s.Color))
	}
	
	// Note: We skip diagonal corner fills for SVG since we don't have diagonal neighbor info
	// The standard implementation checks for diagonal neighbors to fill corners
	// For now, the center circle will cover most visual gaps
	
	// Draw center circle last (on top)
	sb.WriteString(fmt.Sprintf(`<circle cx="%.3f" cy="%.3f" r="%.3f" fill="%s" />`,
		cx, cy, r, s.Color))
	
	return sb.String()
}
