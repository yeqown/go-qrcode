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
