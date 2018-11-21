package version

import (
	"github.com/yeqown/go-qrcode/mask"
	"github.com/yeqown/go-qrcode/matrix"
)

// Version1 ... with size
type Version1 struct {
	vmask        *mask.Mask
	mat          *matrix.Matrix
	recoverLevel RecoveryLevel
}

// Version func return (1-40)
func (v1 *Version1) Version() int {
	return 1
}

// Dimension return v1.size
func (v1 *Version1) Dimension() int {
	return v1.Version()*4 + 17
}

// Init func
func (v1 *Version1) Init(lv RecoveryLevel) error {
	v1.recoverLevel = lv
	v1.mat = matrix.NewMatrix(v1.Dimension(), v1.Dimension())
	v1.vmask = nil
	return nil
}

// Mask func
func (v1 *Version1) Mask() mask.Mask {
	return *v1.vmask
}

// Data func
func (v1 *Version1) Data() *matrix.Matrix {
	return v1.mat
}

// Level ...
func (v1 *Version1) Level() RecoveryLevel {
	return v1.recoverLevel
}

// Cap ... return [num, alpha, byte, jp]
func (v1 *Version1) Cap() [4]int {
	switch v1.recoverLevel {
	case L:
		return [4]int{41, 25, 17, 10}
	case M:
		return [4]int{34, 20, 14, 8}
	case Q:
		return [4]int{27, 16, 11, 7}
	case H:
		return [4]int{17, 10, 7, 4}
	default:
		panic("could not match the recovery level")
	}
}

// ErrorRecoverCap ...
func (v1 *Version1) ErrorRecoverCap() {

}
