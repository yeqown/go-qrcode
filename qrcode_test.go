package qrcode

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	qrc, err := New("cost 3+ days to read docs and coding")
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	t.Logf("analyzed version is: %d, len of data: %d, encMode: %v",
		qrc.v, len(qrc.rawData), qrc.mode)

	if qrc.mode != EncModeByte {
		t.Error("could analyze error with encode type")
		t.Fail()
	}

	// save file
	if err := qrc.Save("./testdata/qrtest.jpeg"); err != nil {
		t.Errorf("could not save image: %v", err)
		t.Fail()
	}

	// check file existed
	_, err = os.Stat("./testdata/qrtest.jpeg")
	if err != nil {
		t.Errorf("could not find image file: %v", err)
		t.Fail()
	}
}

func TestNewWithSpecV(t *testing.T) {
	qrc, err := NewWithSpecV("TestNewWithSpecV", 4, ErrorCorrectionQuart)
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	// save file
	if err := qrc.Save("./testdata/qrtest_spec.jpeg"); err != nil {
		t.Errorf("could not save image: %v", err)
		t.Fail()
	}

	// check file existed
	_, err = os.Stat("./testdata/qrtest_spec.jpeg")
	if err != nil {
		t.Errorf("could not find image file: %v", err)
		t.Fail()
	}
}

func TestNewWithConfig(t *testing.T) {

	encOpts := DefaultConfig()
	encOpts.EcLevel = ErrorCorrectionLow
	encOpts.EncMode = EncModeNumeric

	qrc, err := NewWithConfig("1234567", encOpts, WithQRWidth(20))
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	// save file
	if err := qrc.Save("./testdata/qrtest_spec.jpeg"); err != nil {
		t.Errorf("could not save image: %v", err)
		t.Fail()
	}

	// check file existed
	_, err = os.Stat("./testdata/qrtest_spec.jpeg")
	if err != nil {
		t.Errorf("could not find image file: %v", err)
		t.Fail()
	}
}

// Test_NewWithConfig_UnmatchedEncodeMode NewWithConfig will panic while encMode is
// not matched to Config.EncMode, for example:
// cfg.EncMode is EncModeAlphanumeric but source text is bytes encoding.
func Test_NewWithConfig_UnmatchedEncodeMode(t *testing.T) {
	cfg := DefaultConfig()
	cfg.EncMode = EncModeAlphanumeric

	assert.Panics(t, func() {
		_, err := NewWithConfig("abcs", cfg, WithQRWidth(20))
		if err != nil {
			t.Errorf("could not generate QRCode: %v", err)
			t.Fail()
		}
	})
}

func Test_New_WithOutputOption_bg_fg_width(t *testing.T) {
	qrc, err := New("Test_New_WithOutputOption_bg_fg_width",
		WithBgColorRGBHex("#b8de6f"),
		WithFgColorRGBHex("#01c5c4"),
		WithQRWidth(20),
	)
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	// save file
	if err := qrc.Save("./testdata/qrtest_fg_bg.jpeg"); err != nil {
		t.Errorf("could not save image: %v", err)
		t.Fail()
	}
}

func Test_New_WithOutputOption_Logo(t *testing.T) {
	SetDebugMode()

	qrc, err := New("Test_New_WithOutputOption_Logo",
		WithBgColorRGBHex("#b8de6f"),
		WithFgColorRGBHex("#f1e189"),
		WithLogoImageFilePNG("./testdata/logo.png"), // png required
	)
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	// save file
	if err := qrc.Save("./testdata/qrtest_icon.jpeg"); err != nil {
		t.Errorf("could not save image: %v", err)
		t.Fail()
	}
}

func Test_New_WithOutputOption_Shape(t *testing.T) {
	qrc, err := New("Test_New_WithOutputOption_Shape",
		WithBgColorRGBHex("#b8de6f"),
		WithFgColorRGBHex("#f1e189"),
		WithCircleShape(),
	)
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	// save file
	if err := qrc.Save("./testdata/qrtest_circle.jpeg"); err != nil {
		t.Errorf("could not save image: %v", err)
		t.Fail()
	}
}

func Test_New_WithBorderWidth(t *testing.T) {
	qrc, err := New("Test_New_WithOutputOption_Shape",
		WithBorderWidth(10, 20, 30, 40),
	)
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	// save file
	if err = qrc.Save("./testdata/qrtest_border_width.jpeg"); err != nil {
		t.Errorf("could not save image: %v", err)
		t.Fail()
	}
}

// Test_Issue40
// https://github.com/yeqown/go-qrcode/issues/40
func Test_Issue40(t *testing.T) {
	qrc, err := New("https://baidu.com/")
	if err != nil {
		panic(err)
	}

	err = qrc.Save("./testdata/issue40_1.png")
	if err != nil {
		log.Println(err)
	}
	err = qrc.Save("./testdata/issue40_2.png")
	if err != nil {
		log.Println(err)
	}
}
