package qrcode

import "github.com/yeqown/go-qrcode/v2/matrix"

type Writer interface {
	Write(mat matrix.Matrix) error
}

var _ Writer = (*nonWriter)(nil)

type nonWriter struct{}

func (n nonWriter) Write(mat matrix.Matrix) error {
	return nil
}
