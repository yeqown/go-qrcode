package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

// _defaultOutputOption default output image background color and etc options
var _defaultOutputOption = &outputImageOptions{
	bgColor: color.White, // white
	qrColor: color.Black, // black
	logo:    nil,
	qrWidth: 20,
}

// outputImageOptions to output QR code image
type outputImageOptions struct {
	// bgColors
	bgColor color.Color

	// qrColor
	qrColor color.Color

	// logo this icon image would be put the center of QR Code image
	// NOTE: logo only should has 1/5 size of QRCode image
	logo image.Image

	// qrWidth width of each qr block
	qrWidth int
}

func (oo *outputImageOptions) backgroundColor() color.Color {
	if oo == nil || oo.bgColor == nil {
		return color.White
	}

	return oo.bgColor
}

func (oo *outputImageOptions) foregroundColor() color.Color {
	if oo == nil || oo.qrColor == nil {
		return color.Black
	}

	return oo.qrColor
}

func (oo *outputImageOptions) logoImage() image.Image {
	if oo == nil || oo.logo == nil {
		return nil
	}

	return oo.logo
}

func (oo *outputImageOptions) qrBlockWidth() int {
	if oo == nil || (oo.qrWidth <= 0 || oo.qrWidth > 255) {
		return 20
	}

	return oo.qrWidth
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
		oo.bgColor = c
	})
}

// WithBgColorRGBHex background color
func WithBgColorRGBHex(hex string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.bgColor = hexToRGBA(hex)
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

// WithLogoImage image should only has 1/5 width of QRCode at most
func WithLogoImage(img image.Image) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.logo = img
	})
}

// WithLogoImageFileJPEG load image from file, jpeg is required.
// image should only has 1/5 width of QRCode at most
func WithLogoImageFileJPEG(f string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		fd, err := os.Open(f)
		if err != nil {
			fmt.Printf("could not open file(%s), error=%v\n", f, err)
			return
		}

		img, err := jpeg.Decode(fd)
		if err != nil {
			fmt.Printf("could not open file(%s), error=%v\n", f, err)
			return
		}

		oo.logo = img
	})
}

// WithLogoImageFilePNG load image from file, PNG is required.
// image should only has 1/5 width of QRCode at most
func WithLogoImageFilePNG(f string) ImageOption {
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

		oo.logo = img
	})
}

// WithQRWidth specify width of each qr block
func WithQRWidth(width uint8) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrWidth = int(width)
	})
}
