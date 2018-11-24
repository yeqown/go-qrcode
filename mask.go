package qrcode

import "github.com/yeqown/go-qrcode/matrix"

// MaskPatternModulo ...
// Mask Pattern ref to: https://www.thonky.com/qr-code-tutorial/mask-patterns
type MaskPatternModulo uint32

const (
	// Modulo0 (x+y) mod 2 == 0
	Modulo0 MaskPatternModulo = iota
	// Modulo1 (x) mod 2 == 0
	Modulo1
	// Modulo2 (y) mod 3 == 0
	Modulo2
	// Modulo3 (x+y) mod 3 == 0
	Modulo3
	// Modulo4 (floor (x/ 2) + floor (y/ 3) mod 2 == 0
	Modulo4
	// Modulo5 (x * y) mod 2) + (x * y) mod 3) == 0
	Modulo5
	// Modulo6 (x * y) mod 2) + (x * y) mod 3) mod 2 == 0
	Modulo6
	// Modulo7 (x + y) mod 2) + (x * y) mod 3) mod 2 == 0
	Modulo7
)

// CalculateScore 计算惩罚得分 ...
func CalculateScore(mat *matrix.Matrix) int {
	return 0
}

// Mask ...
type Mask struct {
	mat   *matrix.Matrix    // matrix
	mode  MaskPatternModulo // mode
	score int               // score 惩罚得分，分值越低说明越符合条件
}

// NewMask ...
func NewMask(m *matrix.Matrix, mode MaskPatternModulo) *Mask {
	mask := &Mask{
		mat:  m.Copy(),
		mode: mode,
	}
	mask.init()
	return mask
}

// MaskPatternFunc ...
type MaskPatternFunc func(int, int) bool

// init generate maks by mode
func (m *Mask) init() {
	var f MaskPatternFunc
	switch m.mode {
	case Modulo0:
		f = modulo0Func
	case Modulo1:
		f = modulo1Func
	case Modulo2:
		f = modulo2Func
	case Modulo3:
		f = modulo3Func
	case Modulo4:
		f = modulo4Func
	case Modulo5:
		f = modulo5Func
	case Modulo6:
		f = modulo6Func
	case Modulo7:
		f = modulo7Func
	}

	m.mat.Iter(matrix.ROW, func(x, y int, s matrix.State) {
		// skip the function modules
		if state, _ := m.mat.Get(x, y); state != matrix.StateInit {
			m.mat.Set(x, y, matrix.StateInit)
			return
		}
		if f(x, y) {
			m.mat.Set(x, y, matrix.StateTrue)
		} else {
			m.mat.Set(x, y, matrix.StateFalse)
		}
	})
}

// modulo0Func for maskPattern function
// Modulo0 (x+y) mod 2 == 0
func modulo0Func(x, y int) bool {
	if (x+y)%2 == 0 {
		return true
	}
	return false
}

// modulo1Func for maskPattern function
// Modulo1 (y) mod 2 == 0
func modulo1Func(x, y int) bool {
	if y%2 == 0 {
		return true
	}
	return false
}

// modulo2Func for maskPattern function
// Modulo2 (x) mod 3 == 0
func modulo2Func(x, y int) bool {
	if x%3 == 0 {
		return true
	}
	return false
}

// modulo3Func for maskPattern function
// Modulo3 (x+y) mod 3 == 0
func modulo3Func(x, y int) bool {
	if (x+y)%3 == 0 {
		return true
	}
	return false
}

// modulo4Func for maskPattern function
// Modulo4 (floor (x/ 2) + floor (y/ 3) mod 2 == 0
func modulo4Func(x, y int) bool {
	if (x/3+y/2)%2 == 0 {
		return true
	}
	return false
}

// modulo5Func for maskPattern function
// Modulo5 (x * y) mod 2 + (x * y) mod 3 == 0
func modulo5Func(x, y int) bool {
	if (x*y)%2+(x*y)%3 == 0 {
		return true
	}
	return false
}

// modulo6Func for maskPattern function
// Modulo6 (x * y) mod 2) + (x * y) mod 3) mod 2 == 0
func modulo6Func(x, y int) bool {
	if ((x*y)%2+(x*y)%3)%2 == 0 {
		return true
	}
	return false
}

// modulo7Func for maskPattern function
// Modulo7 (x + y) mod 2) + (x * y) mod 3) mod 2 == 0
func modulo7Func(x, y int) bool {
	if ((x+y)%2+(x*y)%3)%2 == 0 {
		return true
	}
	return false
}
