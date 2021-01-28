package qrcode

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	qrc, err := New("cost 3+ days to read docs and coding")
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
		t.Fail()
	}

	t.Logf("analyzed version is: %d, len of data: %d, encMode: %v",
		qrc.v, len(qrc.rawData), qrc.mode)

	if qrc.mode != encModeByte {
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
	qrc, err := NewWithSpecV("TestNewWithSpecV", 4, Quart)
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
	qrc, err := New("Test_New_WithOutputOption_Logo",
		WithBgColorRGBHex("#b8de6f"),
		WithFgColorRGBHex("#f1e189"),
		//WithCircleShape(),
		//WithLogoImageFilePNG("./testdata/logo.png"), // png required
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
