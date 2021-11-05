package qrcode

import (
	"github.com/yeqown/go-qrcode/matrix"
)

// maskPatternModulo ...
// mask Pattern ref to: https://www.thonky.com/qr-code-tutorial/mask-patterns
type maskPatternModulo uint32

const (
	// modulo0 (x+y) mod 2 == 0
	modulo0 maskPatternModulo = iota
	// modulo1 (x) mod 2 == 0
	modulo1
	// modulo2 (y) mod 3 == 0
	modulo2
	// modulo3 (x+y) mod 3 == 0
	modulo3
	// modulo4 (floor (x/ 2) + floor (y/ 3) mod 2 == 0
	modulo4
	// modulo5 (x * y) mod 2) + (x * y) mod 3) == 0
	modulo5
	// modulo6 (x * y) mod 2) + (x * y) mod 3) mod 2 == 0
	modulo6
	// modulo7 (x + y) mod 2) + (x * y) mod 3) mod 2 == 0
	modulo7
)

var (
	// 1011101 0000
	statePattern1 = []matrix.State{matrix.StateTrue, matrix.StateFalse, matrix.StateTrue, matrix.StateTrue, matrix.StateTrue, matrix.StateFalse, matrix.StateTrue,
		matrix.StateFalse, matrix.StateFalse, matrix.StateFalse, matrix.StateFalse}
	// 0000 1011101
	statePattern2 = []matrix.State{matrix.StateFalse, matrix.StateFalse, matrix.StateFalse, matrix.StateFalse,
		matrix.StateTrue, matrix.StateFalse, matrix.StateTrue, matrix.StateTrue, matrix.StateTrue, matrix.StateFalse, matrix.StateTrue}
)

// calculateScore calculate the maskScore of masking result ...
func calculateScore(mat *matrix.Matrix) int {
	debugLogf("calculate maskScore starting")
	score1 := rule1(mat.Copy())
	score2 := rule2(mat.Copy())
	score3 := rule3(mat.Copy())
	score4 := rule4(mat.Copy())

	debugLogf("maskScore: %d", score1+score2+score3+score4)
	return score1 + score2 + score3 + score4
}

// 第一条规则为一行（或列）中的每组五个或更多相同颜色的模块提供QR代码。
func rule1(mat *matrix.Matrix) int {
	// Row socre
	var (
		score          int
		rowCurState    matrix.State
		rowCurColorCnt int

		colCurState    matrix.State
		colCurColorCnt int
	)

	mat.Iterate(matrix.ROW, func(x, y int, value matrix.State) {
		if x == 0 {
			rowCurColorCnt = 0
			rowCurState = value
			return
		}

		if value == rowCurState {
			rowCurColorCnt++
		} else {
			rowCurState = value
		}

		if rowCurColorCnt == 5 {
			score += 3
		} else if rowCurColorCnt > 5 {
			score++
		}
	})

	// column
	mat.Iterate(matrix.COLUMN, func(x, y int, value matrix.State) {
		if x == 0 {
			colCurColorCnt = 0
			colCurState = value
			return
		}

		if value == colCurState {
			colCurColorCnt++
		} else {
			colCurState = value
		}

		if colCurColorCnt == 5 {
			score += 3
		} else if colCurColorCnt > 5 {
			score++
		}
	})
	return score
}

// 第二个规则给出了QR码对矩阵中相同颜色模块的每个2x2区域的惩罚。
func rule2(mat *matrix.Matrix) int {
	var (
		score          int
		s0, s1, s2, s3 matrix.State
	)
	for x := 0; x < mat.Width()-1; x++ {
		for y := 0; y < mat.Height()-1; y++ {
			s0, _ = mat.Get(x, y)
			s1, _ = mat.Get(x+1, y)
			s2, _ = mat.Get(x, y+1)
			s3, _ = mat.Get(x+1, y+1)

			if s0 == s1 && s2 == s3 && s1 == s2 {
				score += 3
			}
		}
	}

	return score
}

// 如果存在看起来类似于取景器模式的模式，则第三规则给QR码一个大的惩罚
// dark-light-dark-dark-dark-light-dark
// 1011101 0000 or 0000 1011101
func rule3(mat *matrix.Matrix) int {
	var (
		score      int
		stateSlice []matrix.State
	)

	for y := 0; y < mat.Height(); y++ {
		for x := 0; x < mat.Width()-11; x++ {
			for i := 0; i < 11; i++ {
				s, _ := mat.Get(x+i, y)
				stateSlice = append(stateSlice, s)
			}
			if matrix.StateSliceMatched(statePattern1, stateSlice) {
				score += 40
			}
			if matrix.StateSliceMatched(statePattern2, stateSlice) {
				score += 40
			}
		}
	}

	for x := 0; x < mat.Width(); x++ {
		for y := 0; y < mat.Height()-11; y++ {
			// stateSlice =
			for i := 0; i < 11; i++ {
				s, _ := mat.Get(x, y+i)
				stateSlice = append(stateSlice, s)
			}
			if matrix.StateSliceMatched(statePattern1, stateSlice) {
				score += 40
			}
			if matrix.StateSliceMatched(statePattern2, stateSlice) {
				score += 40
			}
		}
	}

	return score
}

// 如果超过一半的模块是暗的或轻的，则第四规则给QR码一个惩罚，对较大的差异有较大的惩罚
func rule4(mat *matrix.Matrix) int {
	var (
		totalCnt             = mat.Width() * mat.Height()
		darkCnt, darkPercent int
	)
	mat.Iterate(matrix.ROW, func(x, y int, s matrix.State) {
		if s == matrix.StateTrue {
			darkCnt++
		}
	})
	darkPercent = (darkCnt * 100) / totalCnt
	x := 0
	if darkPercent%5 == 0 {
		x = 1
	}
	last5Times := abs(((darkPercent/5)-x)*5 - 50)
	next5Times := abs(((darkPercent/5)+1)*5 - 50)

	// get the min maskScore
	if last5Times > next5Times {
		// scoreC <- next5Times / 5 * 10
		return next5Times * 2
	} else {
		return last5Times * 2
	}

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type mask struct {
	mat      *matrix.Matrix    // matrix
	mode     maskPatternModulo // mode
	moduloFn moduloFunc        // moduloFn masking function
}

// newMask ...
func newMask(mat *matrix.Matrix, mode maskPatternModulo) *mask {
	m := &mask{
		mat:      mat.Copy(),
		mode:     mode,
		moduloFn: getModuloFunc(mode),
	}
	m.masking()

	return m
}

// moduloFunc to define what's modulo func
type moduloFunc func(int, int) bool

func getModuloFunc(mode maskPatternModulo) (f moduloFunc) {
	f = nil
	switch mode {
	case modulo0:
		f = modulo0Func
	case modulo1:
		f = modulo1Func
	case modulo2:
		f = modulo2Func
	case modulo3:
		f = modulo3Func
	case modulo4:
		f = modulo4Func
	case modulo5:
		f = modulo5Func
	case modulo6:
		f = modulo6Func
	case modulo7:
		f = modulo7Func
	}

	return
}

// init generate maks by mode
func (m *mask) masking() {
	moduloFn := m.moduloFn
	if moduloFn == nil {
		panic("impossible panic, contact maintainer plz")
	}

	m.mat.Iterate(matrix.ROW, func(x, y int, s matrix.State) {
		// skip the function modules
		if state, _ := m.mat.Get(x, y); state != matrix.StateInit {
			_ = m.mat.Set(x, y, matrix.StateInit)
			return
		}
		if moduloFn(x, y) {
			_ = m.mat.Set(x, y, matrix.StateTrue)
		} else {
			_ = m.mat.Set(x, y, matrix.StateFalse)
		}
	})
}

// modulo0Func for maskPattern function
// modulo0 (x+y) mod 2 == 0
func modulo0Func(x, y int) bool {
	return (x+y)%2 == 0
}

// modulo1Func for maskPattern function
// modulo1 (y) mod 2 == 0
func modulo1Func(x, y int) bool {
	return y%2 == 0
}

// modulo2Func for maskPattern function
// modulo2 (x) mod 3 == 0
func modulo2Func(x, y int) bool {
	return x%3 == 0
}

// modulo3Func for maskPattern function
// modulo3 (x+y) mod 3 == 0
func modulo3Func(x, y int) bool {
	return (x+y)%3 == 0
}

// modulo4Func for maskPattern function
// modulo4 (floor (x/ 2) + floor (y/ 3) mod 2 == 0
func modulo4Func(x, y int) bool {
	return (x/3+y/2)%2 == 0
}

// modulo5Func for maskPattern function
// modulo5 (x * y) mod 2 + (x * y) mod 3 == 0
func modulo5Func(x, y int) bool {
	return (x*y)%2+(x*y)%3 == 0
}

// modulo6Func for maskPattern function
// modulo6 (x * y) mod 2) + (x * y) mod 3) mod 2 == 0
func modulo6Func(x, y int) bool {
	return ((x*y)%2+(x*y)%3)%2 == 0
}

// modulo7Func for maskPattern function
// modulo7 (x + y) mod 2) + (x * y) mod 3) mod 2 == 0
func modulo7Func(x, y int) bool {
	return ((x+y)%2+(x*y)%3)%2 == 0
}
