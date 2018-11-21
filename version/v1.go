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

// Size return v1.size
func (v1 *Version1) Size() int {
	return 1*4 + 17
}

// Init func
func (v1 *Version1) Init(lv RecoveryLevel) error {
	v1.recoverLevel = lv
	v1.mat = matrix.NewMatrix(v1.Size(), v1.Size())
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
