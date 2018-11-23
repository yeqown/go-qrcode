package qrcode

import "github.com/yeqown/go-qrcode/matrix"

// MaskPatterModulo ...
// Mask Pattern ref to: https://www.thonky.com/qr-code-tutorial/mask-patterns
type MaskPatterModulo uint32

const (
	// 0 (row+column) mod 2 == 0
	modulo0 MaskPatterModulo = iota
	// 1 (row) mod 2 == 0
	modulo1
	// 2 (column) mod 3 == 0
	modulo2
	// 3 (row+column) mod 3 == 0
	modulo3
	// 4 (floor (row/ 2) + floor (column/ 3) mod 2 == 0
	modulo4
	// 5 (row * column) mod 2) + (row * column) mod 3) == 0
	modulo5
	// 6 (row * column) mod 2) + (row * column) mod 3) mod 2 == 0
	modulo6
	// 7 (row + column) mod 2) + (row * column) mod 3) mod 2 == 0
	modulo7
)

// Mask ...
type Mask struct {
	mat *matrix.Matrix
}

// New Mask ...
func New(w, h int) *Mask {
	return &Mask{
		mat: matrix.NewMatrix(w, h),
	}
}
