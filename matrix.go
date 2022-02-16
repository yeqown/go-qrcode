package qrcode

import (
	"errors"
	"fmt"
)

var (
	// ErrorOutRangeOfW x out of range of Width
	ErrorOutRangeOfW = errors.New("out of range of width")

	// ErrorOutRangeOfH y out of range of Height
	ErrorOutRangeOfH = errors.New("out of range of height")
)

// newMatrix generate a matrix with map[][]qrbool
func newMatrix(width, height int) *Matrix {
	mat := make([][]qrvalue, width)
	for w := 0; w < width; w++ {
		mat[w] = make([]qrvalue, height)
	}

	m := &Matrix{
		mat:    mat,
		width:  width,
		height: height,
	}

	m.init()
	return m
}

// Matrix is a matrix data type
// width:3 height: 4 for [3][4]int
type Matrix struct {
	mat    [][]qrvalue
	width  int
	height int
}

// do some init work
func (m *Matrix) init() {
	for w := 0; w < m.width; w++ {
		for h := 0; h < m.height; h++ {
			m.mat[w][h] = QRValue_INIT_V0
		}
	}
}

// print to stdout
func (m *Matrix) print() {
	m.Iterate(IterDirection_ROW, func(x, y int, s qrvalue) {
		fmt.Printf("%s ", s)
		if (x + 1) == m.width {
			fmt.Println()
		}
	})
}

// Copy matrix into a new Matrix
func (m *Matrix) Copy() *Matrix {
	mat2 := make([][]qrvalue, m.width)
	for w := 0; w < m.width; w++ {
		mat2[w] = make([]qrvalue, m.height)
		copy(mat2[w], m.mat[w])
	}

	m2 := &Matrix{
		width:  m.width,
		height: m.height,
		mat:    mat2,
	}

	return m2
}

// Width ... width
func (m *Matrix) Width() int {
	return m.width
}

// Height ... height
func (m *Matrix) Height() int {
	return m.height
}

// set [w][h] as true
func (m *Matrix) set(w, h int, c qrvalue) error {
	if w >= m.width || w < 0 {
		return ErrorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return ErrorOutRangeOfH
	}
	m.mat[w][h] = c
	return nil
}

// at state qrbool from matrix with position {x, y}
func (m *Matrix) at(w, h int) (qrvalue, error) {
	if w >= m.width || w < 0 {
		return QRValue_INIT_V0, ErrorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return QRValue_INIT_V0, ErrorOutRangeOfH
	}
	return m.mat[w][h], nil
}

// iterDirection scan matrix direction
type iterDirection uint8

const (
	// IterDirection_ROW for row first
	IterDirection_ROW iterDirection = iota + 1

	// IterDirection_COLUMN for column first
	IterDirection_COLUMN
)

// Iterate the Matrix with loop direction IterDirection_ROW major or IterDirection_COLUMN major.
// IterDirection_COLUMN is recommended.
func (m *Matrix) Iterate(dir iterDirection, f func(int, int, qrvalue)) {
	// row direction first
	if dir == IterDirection_ROW {
		for h := 0; h < m.height; h++ {
			for w := 0; w < m.width; w++ {
				f(w, h, m.mat[w][h])
			}
		}
		return
	}

	// column direction first
	if dir == IterDirection_COLUMN {
		for w := 0; w < m.width; w++ {
			for h := 0; h < m.height; h++ {
				f(w, h, m.mat[w][h])
			}
		}
		return
	}
}

// Row return a row of matrix, cur should be y dimension.
func (m *Matrix) Row(cur int) []qrvalue {
	if cur >= m.height || cur < 0 {
		return nil
	}

	col := make([]qrvalue, m.height)
	for w := 0; w < m.width; w++ {
		col[w] = m.mat[w][cur]
	}
	return col
}

// Col return a slice of column, cur should be x dimension.
func (m *Matrix) Col(cur int) []qrvalue {
	if cur >= m.width || cur < 0 {
		return nil
	}

	return m.mat[cur]
}
