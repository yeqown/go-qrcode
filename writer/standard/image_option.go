package standard

import (
	"github.com/yeqown/go-qrcode/v2"
	"image"
)

type ImageOption interface {
	apply(o *outputImageOptions)
}

// defaultOutputImageOption default output image background color and etc options
func defaultOutputImageOption() *outputImageOptions {
	return &outputImageOptions{
		bgColor:       color_WHITE_rgba, // rgba white
		bgTransparent: false,            // not transparent
		qrColor:       color_BLACK_rgba, // rgba black
		compressed:    false,            // disable compression mode
		logo:          nil,              //
		qrWidth:       20,               //
		shape:         _shapeRectangle,  //
		imageEncoder:  jpegEncoder{},
		borderWidths:  [4]int{_defaultPadding, _defaultPadding, _defaultPadding, _defaultPadding},
	}
}

// outputImageOptions to output QR code image
type outputImageOptions struct {
	// bgColor is the background color of the QR code image, qrColor is
	// the foreground color of the QR code.
	bgColor, qrColor Color

	// bgTransparent only affects on PNG_FORMAT
	bgTransparent bool
	// compressed represents output image should be generated in the minimum size. only
	// support gray color in compression mode, so APIs would take no any effect to output, such
	// as WithBgColor, WithFgColor etc.
	compressed bool

	// logo this icon image would be put the center of QR Code image
	// NOTE: logo only should have 1/5 size of QRCode image
	logo image.Image

	// qrWidth width of each qr block
	qrWidth int

	// shape means how to draw the shape of each cell.
	shape IShape

	// imageEncoder specify which file format would be encoded the QR image.
	imageEncoder ImageEncoder

	// borderWidths indicates the border width of the output image. the order is
	// top, right, bottom, left same as the WithBorder
	borderWidths [4]int

	// halftoneImg is the halftone image for the output image.
	halftoneImg image.Image
}

func (oo *outputImageOptions) backgroundColor() Color {
	if oo == nil {
		return color_WHITE
	}

	return oo.bgColor
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

func (oo *outputImageOptions) getShape() IShape {
	if oo == nil || oo.shape == nil {
		return _shapeRectangle
	}

	return oo.shape
}

// preCalculateAttribute this function must reference to draw function.
func (oo *outputImageOptions) preCalculateAttribute(dimension int) *Attribute {
	if oo == nil {
		return nil
	}

	top, right, bottom, left := oo.borderWidths[0], oo.borderWidths[1], oo.borderWidths[2], oo.borderWidths[3]
	return &Attribute{
		W:          dimension*oo.qrBlockWidth() + right + left,
		H:          dimension*oo.qrBlockWidth() + top + bottom,
		Borders:    oo.borderWidths,
		BlockWidth: oo.qrBlockWidth(),
	}
}

// translateQrColor get color.RGBA by value State, if not found, return outputImageOptions.qrColor.
// NOTE: this function decides the state should use qrColor or bgColor.
func (oo *outputImageOptions) translateQrColor(v qrcode.QRValue) (c Color) {
	// TODO(@yeqown): use _STATE_MAPPING to replace this function while in debug mode
	// or some special flag.
	if v.IsSet() {
		return oo.qrColor
	}

	return oo.bgColor
}
