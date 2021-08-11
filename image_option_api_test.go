package qrcode

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeqown/go-qrcode/matrix"
)

func Test_WithBuiltinImageEncoder(t *testing.T) {
	oo := defaultOutputImageOption()

	assert.IsType(t, jpegEncoder{}, oo.imageEncoder)
	WithBuiltinImageEncoder(JPEG_FORMAT).apply(oo)
	assert.IsType(t, jpegEncoder{}, oo.imageEncoder)
	WithBuiltinImageEncoder(PNG_FORMAT).apply(oo)
	assert.IsType(t, pngEncoder{}, oo.imageEncoder)
}

func Test_WithCustomImageEncoder(t *testing.T) {
	oo := defaultOutputImageOption()

	assert.IsType(t, jpegEncoder{}, oo.imageEncoder)
	WithCustomImageEncoder(nil).apply(oo)
	assert.IsType(t, jpegEncoder{}, oo.imageEncoder)
}

func Test_BgColor_FgColor(t *testing.T) {
	oo := defaultOutputImageOption()

	// check
	assert.Equal(t, color.White, oo.bgColor)
	assert.Equal(t, color.Black, oo.qrColor)

	// apply color
	WithBgColor(color.Black).apply(oo)
	assert.Equal(t, color.Black, oo.bgColor)
	assert.Equal(t, color.Black, oo.qrColor)

	// apply color
	WithBgColor(color.White).apply(oo)
	WithFgColor(color.White).apply(oo)
	assert.Equal(t, color.White, oo.bgColor)
	assert.Equal(t, color.White, oo.qrColor)
	assert.Equal(t, color.White, _stateToRGBA[matrix.StateFinder])

	WithFgColor(color.Black).apply(oo)
	assert.Equal(t, color.White, oo.bgColor)
	assert.Equal(t, color.Black, oo.qrColor)
	assert.Equal(t, color.Black, _stateToRGBA[matrix.StateFinder])
}

func Test_defaultOutputOption(t *testing.T) {
	oo := defaultOutputImageOption()

	// Apply
	rgba := color.RGBA{
		R: 123,
		G: 123,
		B: 123,
		A: 123,
	}
	WithBgColor(rgba).apply(oo)
	// assert
	assert.Equal(t, rgba, oo.bgColor)

	// check default
	oo2 := defaultOutputImageOption()
	assert.NotEqual(t, oo2.bgColor, oo.bgColor)
}

func Test_WithBorderWidth(t *testing.T) {
	oo := defaultOutputImageOption()

	// zero parameter
	WithBorderWidth().apply(oo)
	assert.Equal(t, [4]int{_defaultPadding, _defaultPadding, _defaultPadding, _defaultPadding}, oo.borderWidths)

	// one parameter
	WithBorderWidth(1).apply(oo)
	assert.Equal(t, [4]int{1, 1, 1, 1}, oo.borderWidths)

	// two parameters
	WithBorderWidth(1, 2).apply(oo)
	assert.Equal(t, [4]int{1, 2, 1, 2}, oo.borderWidths)

	// three parameters
	WithBorderWidth(1, 2, 3).apply(oo)
	assert.Equal(t, [4]int{1, 2, 1, 2}, oo.borderWidths)

	// four parameters
	WithBorderWidth(1, 2, 3, 4).apply(oo)
	assert.Equal(t, [4]int{1, 2, 3, 4}, oo.borderWidths)
}
