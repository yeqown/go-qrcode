package matrix

import (
	"errors"
	"fmt"
	"image/color"
)

// DIRECTION scan matrix driection
type DIRECTION uint

// State value of matrix map[][]
type State uint16

const (
	// ROW for row first
	ROW DIRECTION = 1
	// COLUMN for column first
	COLUMN DIRECTION = 2

	// StateFalse 0xffff FLASE
	StateFalse State = 0xffff

	// ZERO 0x0 FALSE
	ZERO State = 0xeeee
	// StateTrue 0x0 TRUE
	StateTrue State = 0x0

	// StateInit 0x9999 use for initial state
	StateInit State = 0x9999

	// StateVersion 0x4444
	StateVersion State = 0x4444

	// StateFormat 0x1234 for presisted state
	StateFormat State = 0x7777

	// BORDER ... gray color
	BORDER State = 0x3333
)

// LoadGray16 load color by value State
func LoadGray16(v State) color.Gray16 {
	switch v {
	case StateFalse:
		return colorMap[StateFalse]
	case StateTrue:
		return colorMap[StateTrue]
	case StateInit:
		return colorMap[StateInit]
	case StateVersion:
		return colorMap[StateVersion]
	case StateFormat:
		return colorMap[StateFormat]
	}
	return color.Gray16{Y: uint16(v)}
}

var (
	// ErrorOutRangeOfW x out of range of Width
	ErrorOutRangeOfW = errors.New("out of range of width")

	// ErrorOutRangeOfH y out of range of Height
	ErrorOutRangeOfH = errors.New("out of range of height")

	// colorMap state map tp color.Gray16
	colorMap = map[State]color.Gray16{
		StateFalse:   color.Gray16{Y: uint16(StateFalse)},
		StateTrue:    color.Gray16{Y: uint16(StateTrue)},
		StateInit:    color.Gray16{Y: uint16(StateInit)},
		StateVersion: color.Gray16{Y: uint16(StateVersion)},
		StateFormat:  color.Gray16{Y: uint16(StateFormat)},
	}
)

// New generate a matrix with map[][]bool
func New(width, height int) *Matrix {
	mat := make([][]State, width)
	for w := 0; w < width; w++ {
		mat[w] = make([]State, height)
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
	mat    [][]State
	width  int
	height int
}

// do some init work
func (m *Matrix) init() {
	for w := 0; w < m.width; w++ {
		for h := 0; h < m.height; h++ {
			m.mat[w][h] = StateInit
		}
	}
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

// XOR operation ...
// TODO:
func (m *Matrix) XOR() *Matrix {
	return nil
}

// Copy matrix into a new Matrix
func (m *Matrix) Copy() *Matrix {
	newMat := make([][]State, m.width)
	for w := 0; w < m.width; w++ {
		newMat[w] = make([]State, m.height)
		copy(newMat[w], m.mat[w])
	}

	newM := &Matrix{
		width:  m.width,
		height: m.height,
		mat:    newMat,
	}
	return newM
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
func (m *Matrix) Set(w, h int, c State) error {
	if w >= m.width || w < 0 {
		return ErrorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return ErrorOutRangeOfH
	}
	m.mat[w][h] = c
	return nil
}

// Reset [w][h] as false
func (m *Matrix) Reset(w, h int) error {
	return m.Set(w, h, ZERO)
}

// Get ... from mat
func (m *Matrix) Get(w, h int) (State, error) {
	if w >= m.width || w < 0 {
		return ZERO, ErrorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return ZERO, ErrorOutRangeOfH
	}
	return m.mat[w][h], nil
}

// IterFunc ...
type IterFunc func(int, int, State)

// Iter the Matrix
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

// XOR ...
func XOR(s1, s2 State) State {
	if s1 != s2 {
		return StateTrue
	}
	return StateFalse
}
