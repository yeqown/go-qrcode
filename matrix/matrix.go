package matrix

import (
	"errors"
	"fmt"
)

// DIRECTION scan matrix driection
type DIRECTION uint

const (
	// ROW for row first
	ROW DIRECTION = 1
	// COLUMN for column first
	COLUMN DIRECTION = 2
)

var (
	errorOutRangeOfW = errors.New("out of range of width")
	errorOutRangeOfH = errors.New("out of range of height")
)

func NewMatrix(width, height int) *Matrix {

	mat := make([][]bool, width)
	for w := 0; w < width; w++ {
		mat[w] = make([]bool, height)
	}

	m := &Matrix{
		mat:    mat,
		width:  width,
		height: height,
	}

	m.init()
	return m
}

// Matrix ...
// width:3 height: 4 for [3][4]int
//
type Matrix struct {
	mat    [][]bool
	width  int
	height int
}

// do some init work
func (m *Matrix) init() {

}

// Print to stdout
func (m *Matrix) print() {
	fmt.Println("====== matrix =======")
	for h := 0; h < m.height; h++ {
		for w := 0; w < m.width; w++ {
			fmt.Printf("%v ", m.mat[w][h])
		}
		fmt.Println()
	}
	fmt.Println()
}

// Width ... width
func (m *Matrix) Width() int {
	return m.width
}

// Height ... height
func (m *Matrix) Height() int {
	return m.height
}

// Set [w][h] as true
func (m *Matrix) Set(w, h int) error {
	if w >= m.width || w < 0 {
		return errorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return errorOutRangeOfH
	}
	m.mat[w][h] = true
	return nil
}

// Reset [w][h] as false
func (m *Matrix) Reset(w, h int) error {
	if w >= m.width || w < 0 {
		return errorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return errorOutRangeOfH
	}
	m.mat[w][h] = false
	return nil
}

// Get ... from mat
func (m *Matrix) Get(w, h int) (bool, error) {
	if w >= m.width || w < 0 {
		return false, errorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return false, errorOutRangeOfH
	}
	return m.mat[w][h], nil
}

// IterFunc ...
type IterFunc func(int, int, bool)

func (m *Matrix) Iter(dir DIRECTION, f IterFunc) {
	// row first 行优先
	if dir == ROW {
		for h := 0; h < m.height; h++ {
			for w := 0; w < m.width; w++ {
				f(w, h, m.mat[w][h])
			}
		}
		return
	}

	// column first 列优先
	if dir == COLUMN {
		for w := 0; w < m.width; w++ {
			for h := 0; h < m.height; h++ {
				f(w, h, m.mat[w][h])
			}
		}
		return
	}
}
