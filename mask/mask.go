package mask

import "github.com/yeqown/go-qrcode/matrix"

/*
 QRCode mask
 ref to:
*/

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
