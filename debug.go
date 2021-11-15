package qrcode

import (
	"errors"

	"github.com/yeqown/go-qrcode/v2/matrix"
)

// TODO(@yeqown): save matrix in debug mode
func debugDraw(filename string, mat matrix.Matrix) error {
	if _debug {
		return nil
	}

	return errors.New("implement me")
}
