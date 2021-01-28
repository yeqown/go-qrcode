package qrcode

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithBuiltinImageEncoder(t *testing.T) {
	oo := _defaultOutputOption

	assert.IsType(t, oo.imageEncoder, jpegEncoder{})
	WithBuiltinImageEncoder(JPEG_FORMAT).apply(oo)
	assert.IsType(t, oo.imageEncoder, jpegEncoder{})
	WithBuiltinImageEncoder(PNG_FORMAT).apply(oo)
	assert.IsType(t, oo.imageEncoder, pngEncoder{})
}

func Test_WithCustomImageEncoder(t *testing.T) {
	oo := _defaultOutputOption

	assert.IsType(t, oo.imageEncoder, jpegEncoder{})
	WithCustomImageEncoder(nil).apply(oo)
	assert.IsType(t, oo.imageEncoder, jpegEncoder{})
}

func Test_BgColor_FgColor(t *testing.T) {
	oo := _defaultOutputOption

	if oo.bgColor != color.White || oo.qrColor != color.Black {
		t.Error("default value failed")
		t.FailNow()
	}

	WithBgColor(color.Black).apply(oo)
	if oo.bgColor != color.Black || oo.qrColor != color.Black {
		t.Error("value set failed")
		t.FailNow()
	}
	WithBgColor(color.White).apply(oo)
	WithFgColor(color.White).apply(oo)
	if oo.bgColor != color.White || oo.qrColor != color.White {
		t.Error("value set failed")
		t.FailNow()
	}

	WithFgColor(color.Black).apply(oo)
	if oo.bgColor != color.White || oo.qrColor != color.Black {
		t.Error("value set failed")
		t.FailNow()
	}
}
