package qart

import (
	"testing"

	"github.com/yeqown/go-qrcode/v2"
)

func Test_New(t *testing.T) {
	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		t.Fatal(err)
	}

	w, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// save file
	err = qrc.Save(w)
}
