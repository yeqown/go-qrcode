package beautifier

import (
	"github.com/pkg/errors"

	"github.com/yeqown/go-qrcode/v2"
)

// Writer is a QArt output writer that helps to replace some fixed patterns in
// qr matrix to customized images.
type Writer struct {
	// enabledPatterns indicate patterns that would take effect
	// to generate QR image.
	enabledPatterns []patternDesc
}

func New() (qrcode.Writer, error) {
	w := Writer{}
	err := w.preload()
	if err != nil {
		return nil, errors.Wrap(err, "preload failed")
	}

	return w, nil
}

const __OCCUPIED = uint8(0x1a) // 00011011

// Write executes as following steps:
// 1. find all available patterns and record their information, such as: position and size
// 2. now rewrite output image with patterns.
// 3. if one block did not matched in any of the patterns, draw with default image.
func (w Writer) Write(mat qrcode.Matrix) error {
	occupied := mat.Copy()
	// occupied.

	ignore := func(x, y int, s qrcode.QRValue) bool {
		if v := occupied.Get(x, y); v == __OCCUPIED {
			return true
		}

		// for now finder, can't be customized.
		if s.Type() == qrcode.QRType_FINDER {
			return true
		}

		return false
	}

	mat.Iterate(qrcode.IterDirection_COLUMN, func(x, y int, s qrcode.QRValue) {
		// if the position has been occupied, or represents special symbol such as position marks,
		// it can't be handled, skip and handle next position.
		if ignore(x, y, s) {
			return
		}

		for _, pd := range w.enabledPatterns {
			applyPattern(mat, occupied, x, y, pd)
		}
	})

	// TODO(@yeqown): generate image file and save.

	return nil
}

func (w Writer) Close() error {
	// TODO(@yeqown): release pattern resource files.
	return nil
}

func (w Writer) preload() error {
	// TODO(@yeqown): load resource files
	return nil
}

// applyPattern
func applyPattern(mat qrcode.Matrix, occupied *qrcode.Matrix, x, y int, pattern patternDesc) bool {
	return false
}
