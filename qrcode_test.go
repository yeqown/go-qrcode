package qrcode

import (
	"os"
	"testing"
)

func TestQRCOde(t *testing.T) {
	qrc, err := New("http://baidu.com")
	if err != nil {
		t.Errorf("could not generate QRCode: %v", err)
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
