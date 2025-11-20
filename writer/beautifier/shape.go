package beautifier

// Point represents a coordinate in the QR matrix.
type Point struct {
	X, Y int
}


// ShapeType defines the type of a shape.
type ShapeType int

const (
	ShapeTypeBlock ShapeType = iota
	ShapeTypeFinder
	ShapeTypeAlignment
	ShapeTypeTiming
	ShapeTypeLineX
	ShapeTypeLineY
	ShapeTypeSquare
	ShapeTypeL
)

// Shape represents a detected visual structure in the QR code.
type Shape struct {
	Type   ShapeType
	Points []Point
	// Meta can hold additional info like color, Neighbors, etc.
	// For now, we might not need it if Points are sufficient.
	// But for context-aware styles (like Liquid), we might need neighbors.
	// Actually, the Style implementation can check the Matrix if needed, 
	// provided we pass the Matrix to Draw.
}

// Point is already defined in pattern.go, but we might want to move it here or keep it there.
// Since we are refactoring, let's assume pattern.go might be deprecated or merged.
// For now, let's rely on pattern.go's Point if it exists, or redefine if we delete pattern.go.
// We previously had Point in pattern.go.
