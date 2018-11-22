package qrcode

import (
	"testing"
)

// func TestInit() {
// 	load(defaultPathToCfg)
// }

func TestEncodeNum(t *testing.T) {
	enc := Encoder{
		mode:    EncModeNumeric,
		version: loadVersion(1, L),
	}

	b, err := enc.Encode([]byte("12312312"))
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.At(2))
}

func TestEncodeAlphanum(t *testing.T) {
	enc := Encoder{
		mode:    EncModeAlphanumeric,
		version: loadVersion(1, L),
	}

	b, err := enc.Encode([]byte("AKJA*:/"))
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.At(2))
}

func TestEncodeByte(t *testing.T) {
	enc := Encoder{
		mode:    EncModeByte,
		version: loadVersion(1, L),
	}

	b, err := enc.Encode([]byte("helloworld"))
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.At(2))
}
