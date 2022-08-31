package standard

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, color_WHITE, oo.bgColor)
	assert.Equal(t, color_BLACK, oo.qrColor)

	// apply color
	WithBgColor(color_BLACK).apply(oo)
	assert.Equal(t, color_BLACK, oo.bgColor)
	assert.Equal(t, color_BLACK, oo.qrColor)

	// apply color
	WithBgColor(color_WHITE).apply(oo)
	WithFgColor(color_WHITE).apply(oo)
	assert.Equal(t, color_WHITE, oo.bgColor)
	assert.Equal(t, color_WHITE, oo.qrColor)

	WithFgColor(color_BLACK).apply(oo)
	assert.Equal(t, color_WHITE, oo.bgColor)
	assert.Equal(t, color_BLACK, oo.qrColor)
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

func Test_WithBgTransparent(t *testing.T) {
	oo := defaultOutputImageOption()

	// check default
	assert.False(t, oo.bgTransparent)

	// apply
	WithBgTransparent().apply(oo)
	assert.True(t, oo.bgTransparent)
}

func Test_WithLogoSizeMultiplier(t *testing.T) {
	oo := defaultOutputImageOption()

	// check default
	assert.Equal(t, 5, oo.logoSizeMultiplier)

	// apply
	WithLogoSizeMultiplier(2).apply(oo)
	assert.Equal(t, 2, oo.logoSizeMultiplier)
}
