package file

import (
	"errors"
	"os"

	"github.com/yeqown/go-qrcode/v2"
)

const (
	upRune     = 9600 // '▀'
	downRune   = 9604 // '▄'
	upDownRune = 9608 // '█'
	spaceRune  = 32   // ' '
)

var _ qrcode.Writer = (*Writer)(nil)

// Writer implements qrcode.Writer.
type Writer struct {
	out *os.File
}

// Close method to implement qrcode.Writer.
func (a *Writer) Close() error { return nil }

// Write method to implement qrcode.Writer.
func (a *Writer) Write(mat qrcode.Matrix) error {
	if a.out == nil {
		return errors.New("nil file")
	}

	bm := mat.Bitmap()

	output := make([][]string, len(bm)/2+len(bm)%2)
	for t := range output {
		output[t] = make([]string, len(bm[0]))
	}

	for col := range output {
		for row := range output[col] {
			var selectedRune rune = spaceRune

			if bm[col*2][row] {
				selectedRune = upRune
			}

			if col*2+1 < len(bm) {
				if bm[col*2+1][row] && !bm[col*2][row] {
					selectedRune = downRune
				}

				if bm[col*2+1][row] && bm[col*2][row] {
					selectedRune = upDownRune
				}
			}

			output[col][row] = string(selectedRune)
		}
	}

	for _, col := range output {
		for _, row := range col {
			_, err := a.out.WriteString(row)
			if err != nil {
				return err
			}
		}
		_, err := a.out.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func New(f *os.File) *Writer {
	return &Writer{out: f}
}
