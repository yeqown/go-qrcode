package qrcode

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

type formatTyp uint8

const (
	// JPEG_FORMAT as default output file format.
	JPEG_FORMAT formatTyp = iota
	// PNG_FORMAT .
	PNG_FORMAT
	// HEIF_FORMAT High Efficiency Image File Format
	HEIF_FORMAT
)

type ImageEncoder interface {
	// Encode specify which format to encode image into w io.Writer.
	Encode(w io.Writer, img image.Image) error
}

type jpegEncoder struct{}

func (j jpegEncoder) Encode(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, nil)
}

type pngEncoder struct{}

func (j pngEncoder) Encode(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

type heifEncoder struct{}

func (j heifEncoder) Encode(w io.Writer, img image.Image) error {
	panic("Not implemented")
}
