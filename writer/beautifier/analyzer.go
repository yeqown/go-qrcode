package beautifier

import (
	"github.com/yeqown/go-qrcode/v2"
)

// Analyze scans the matrix and returns a list of Shapes.
func Analyze(mat qrcode.Matrix) []Shape {
	var shapes []Shape
	width := mat.Width()
	height := mat.Height()
	visited := make(map[Point]bool)

	// Helper to check if a point is visited
	isVisited := func(x, y int) bool {
		return visited[Point{x, y}]
	}

	// Helper to mark points as visited
	markVisited := func(points []Point) {
		for _, p := range points {
			visited[p] = true
		}
	}

	// 1. Detect Function Patterns (Finder, Alignment, Timing)
	// We iterate and check QRType.
	// Note: The matrix stores QRType for every block.
	// We can group them.

	// Finders (7x7)
	// We know where they are: (0,0), (w-7, 0), (0, w-7).
	finderBases := []Point{{0, 0}, {width - 7, 0}, {0, width - 7}}
	for _, base := range finderBases {
		// Verify it is Finder type (just to be safe)
		if mat.Col(base.X)[base.Y].Type() == qrcode.QRType_FINDER {
			points := make([]Point, 0, 49)
			for i := 0; i < 7; i++ {
				for j := 0; j < 7; j++ {
					points = append(points, Point{base.X + i, base.Y + j})
				}
			}
			shapes = append(shapes, Shape{Type: ShapeTypeFinder, Points: points})
			markVisited(points)
		}
	}

	// Other Function Patterns (Alignment, Timing, Version, Format)
	// We can treat them as generic "Block" or specific types if we want to style them differently.
	// For now, let's group them as "Block" or specific if requested.
	// The user mentioned "Finder" specifically.
	// Let's iterate and find other special types.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if isVisited(x, y) {
				continue
			}
			v := mat.Col(x)[y]
			if !v.IsSet() {
				continue
			}

			t := v.Type()
			if t == qrcode.QRType_TIMING || t == qrcode.QRType_VERSION || t == qrcode.QRType_FORMAT {
				// Treat as single block for now, but with specific type if we want to support it later.
				// Or we can just emit ShapeTypeBlock.
				// Let's use ShapeTypeBlock but we might want to protect them?
				// Actually, if we want to style them, we should use specific types.
				// Let's map them to ShapeTypeAlignment / ShapeTypeTiming etc.
				
				var st ShapeType
				switch t {
				case qrcode.QRType_TIMING:
					st = ShapeTypeTiming
				default:
					st = ShapeTypeBlock // Version/Format treated as Block for now
				}
				
				shapes = append(shapes, Shape{Type: st, Points: []Point{{x, y}}})
				markVisited([]Point{{x, y}})
			}
		}
	}

	// 2. Detect Geometric Patterns in Data
	// We iterate again for remaining unvisited set blocks (which should be Data).
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isVisited(x, y) {
				continue
			}
			if !mat.Col(x)[y].IsSet() {
				continue
			}

			// Try to match shapes
			// Priority: Square -> L -> LineX -> LineY -> Block

			// Square (2x2)
			if x+1 < width && y+1 < height &&
				!isVisited(x+1, y) && !isVisited(x, y+1) && !isVisited(x+1, y+1) &&
				mat.Col(x+1)[y].IsSet() && mat.Col(x)[y+1].IsSet() && mat.Col(x+1)[y+1].IsSet() {
				
				points := []Point{{x, y}, {x + 1, y}, {x, y + 1}, {x + 1, y + 1}}
				shapes = append(shapes, Shape{Type: ShapeTypeSquare, Points: points})
				markVisited(points)
				continue
			}

			// L Shape (3 blocks)
			// Corner at (x,y), Right (x+1,y), Bottom (x,y+1)
			if x+1 < width && y+1 < height &&
				!isVisited(x+1, y) && !isVisited(x, y+1) &&
				mat.Col(x+1)[y].IsSet() && mat.Col(x)[y+1].IsSet() {
				
				points := []Point{{x, y}, {x + 1, y}, {x, y + 1}}
				shapes = append(shapes, Shape{Type: ShapeTypeL, Points: points})
				markVisited(points)
				continue
			}

			// Line X (>=2)
			if x+1 < width && !isVisited(x+1, y) && mat.Col(x+1)[y].IsSet() {
				points := []Point{{x, y}, {x + 1, y}}
				currX := x + 2
				for currX < width && !isVisited(currX, y) && mat.Col(currX)[y].IsSet() {
					points = append(points, Point{currX, y})
					currX++
				}
				shapes = append(shapes, Shape{Type: ShapeTypeLineX, Points: points})
				markVisited(points)
				continue
			}

			// Line Y (>=2)
			if y+1 < height && !isVisited(x, y+1) && mat.Col(x)[y+1].IsSet() {
				points := []Point{{x, y}, {x, y + 1}}
				currY := y + 2
				for currY < height && !isVisited(x, currY) && mat.Col(x)[currY].IsSet() {
					points = append(points, Point{x, currY})
					currY++
				}
				shapes = append(shapes, Shape{Type: ShapeTypeLineY, Points: points})
				markVisited(points)
				continue
			}

			// Single Block
			points := []Point{{x, y}}
			shapes = append(shapes, Shape{Type: ShapeTypeBlock, Points: points})
			markVisited(points)
		}
	}

	return shapes
}
