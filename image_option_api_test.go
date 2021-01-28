package qrcode

import (
	"image/color"
	"testing"
)

func Test_WithOutputFormat(t *testing.T) {
	oo := _defaultOutputOption

	if oo.fileFormat != JPEG_FORMAT {
		t.Error("default value failed")
		t.FailNow()
	}
	WithOutputFormat(JPEG_FORMAT).apply(oo)
	if oo.fileFormat != JPEG_FORMAT {
		t.Error("value set failed")
		t.FailNow()
	}
	WithOutputFormat(PNG_FORMAT).apply(oo)
	if oo.fileFormat != PNG_FORMAT {
		t.Error("value set failed")
		t.FailNow()
	}
	WithOutputFormat(HEIF_FORMAT).apply(oo)
	if oo.fileFormat != HEIF_FORMAT {
		t.Error("value set failed")
		t.FailNow()
	}
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
