package qrcode

import (
	"image"
	"image/color"
)

// _defaultOutputOption default output image background color and etc options
var _defaultOutputOption = &outputImageOptions{
	backgroundColor: color.White, // white
	qrColor:         color.Black, // black
	logoImage:       nil,
}

// outputImageOptions to output QR code image
type outputImageOptions struct {
	// backgroundColor
	backgroundColor color.Color

	// qrColor
	qrColor color.Color

	// logoImage this icon image would be put the center of QR Code image
	logoImage image.Image
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

// WithBgColor .
func WithBgColor(c color.Color) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.backgroundColor = c
	})
}

// WithBgColorRGBHex .
func WithBgColorRGBHex(hex string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.backgroundColor = hexToRGBA(hex)
	})
}

// WithFgColor .
func WithFgColor(c color.Color) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrColor = c
	})
}

// WithFgColorRGBHex "#123123"
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

// WithLogoImageFile .
// TODO: load image from file, PNG is required
func WithLogoImageFile(f string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		var img image.Image
		oo.logoImage = img
	})
}
