package compressed

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	"github.com/yeqown/go-qrcode/v2"
)

type Option struct {
	Padding   int
	BlockSize int
}

// compressedWriter implements issue#69, generating compressed images
// in some special situations, such as, network transferring.
// https://github.com/yeqown/go-qrcode/issues/69
type compressedWriter struct {
	fd io.WriteCloser

	option *Option
}

var (
	backgroundColor = color.Gray{Y: 0xff}
	foregroundColor = color.Gray{Y: 0x00}
)

func New(filename string, opt *Option) (qrcode.Writer, error) {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return compressedWriter{fd: fd, option: opt}, nil
}

func NewWithWriter(writeCloser io.WriteCloser, opt *Option) qrcode.Writer {
	return compressedWriter{fd: writeCloser, option: opt}
}

func (w compressedWriter) Write(mat qrcode.Matrix) error {
	padding := w.option.Padding
	blockWidth := w.option.BlockSize
	width := mat.Width()*blockWidth + 2*padding
	height := width

	img := image.NewPaletted(
		image.Rect(0, 0, width, height),
		[]color.Color{backgroundColor, foregroundColor},
	)
	bgColor := uint8(img.Palette.Index(backgroundColor))
	fgColor := uint8(img.Palette.Index(foregroundColor))

	rectangle := func(x1, y1 int, x2, y2 int, img *image.Paletted, color uint8) {
		for x := x1; x < x2; x++ {
			for y := y1; y < y2; y++ {
				pos := img.PixOffset(x, y)
				img.Pix[pos] = color
			}
		}
	}

	// background
	rectangle(0, 0, width, height, img, bgColor)

	mat.Iterate(qrcode.IterDirection_COLUMN, func(x int, y int, v qrcode.QRValue) {
		sx := x*blockWidth + padding
		sy := y*blockWidth + padding
		es := (x+1)*blockWidth + padding
		ey := (y+1)*blockWidth + padding

		if v.IsSet() {
			rectangle(sx, sy, es, ey, img, fgColor)
		}

		//switch v.IsSet() {
		//case false:
		//	gray = backgroundColor
		//default:
		//	gray = foregroundColor
		//}

	})

	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	return encoder.Encode(w.fd, img)
}

func (w compressedWriter) Close() error {
	return w.fd.Close()
}
