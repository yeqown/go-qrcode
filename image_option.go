package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// _defaultOutputOption default output image background color and etc options
var _defaultOutputOption = &outputImageOptions{
	backgroundColor: color.White, // white
	qrColor:         color.Black, // black
	logoImage:       nil,
	qrWidth:         20,
}

// outputImageOptions to output QR code image
type outputImageOptions struct {
	// backgroundColor
	backgroundColor color.Color

	// qrColor
	qrColor color.Color

	// logoImage this icon image would be put the center of QR Code image
	logoImage image.Image

	// qrWidth width of each qr block
	qrWidth int
}

type ImageOption interface {
	apply(o *outputImageOptions)
}

// funcOption wraps a function that modifies outputImageOptions into an
// implementation of the ImageOption interface.
type funcOption struct {
	f func(oo *outputImageOptions)
}

func (fo *funcOption) apply(oo *outputImageOptions) {
	fo.f(oo)
}

func newFuncDialOption(f func(oo *outputImageOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithBgColor background color
func WithBgColor(c color.Color) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.backgroundColor = c
	})
}

// WithBgColorRGBHex background color
func WithBgColorRGBHex(hex string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.backgroundColor = hexToRGBA(hex)
	})
}

// WithFgColor QR color
func WithFgColor(c color.Color) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrColor = c
	})
}

// WithFgColorRGBHex Hex string to set QR Color
func WithFgColorRGBHex(hex string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrColor = hexToRGBA(hex)
	})
}

// WithLogoImage .
func WithLogoImage(img image.Image) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.logoImage = img
	})
}

// WithLogoImageFile load image from file, PNG is required
func WithLogoImageFile(f string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		fd, err := os.Open(f)
		if err != nil {
			fmt.Printf("could not open file(%s), error=%v\n", f, err)
			return
		}

		img, err := png.Decode(fd)
		if err != nil {
			fmt.Printf("could not open file(%s), error=%v\n", f, err)
			return
		}

		oo.logoImage = img
	})
}

// WithQRWidth specify width of each qr block
func WithQRWidth(width uint8) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrWidth = int(width)
	})
}
