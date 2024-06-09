package qrcode

import (
	"testing"
)

func TestEncodeNum(t *testing.T) {
	enc := encoder{
		ecLv:    ErrorCorrectionLow,
		mode:    EncModeNumeric,
		version: loadVersion(1, ErrorCorrectionLow),
	}

	b, err := enc.Encode("12312312")
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.Len())
}

func TestEncodeAlphanum(t *testing.T) {
	enc := encoder{
		ecLv:    ErrorCorrectionLow,
		mode:    EncModeAlphanumeric,
		version: loadVersion(1, ErrorCorrectionLow),
	}

	b, err := enc.Encode("AKJA*:/")
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.Len())
}

func TestEncodeByte(t *testing.T) {
	enc := encoder{
		ecLv:    ErrorCorrectionQuart,
		mode:    EncModeByte,
		version: loadVersion(5, ErrorCorrectionQuart),
	}

	b, err := enc.Encode("http://baidu.com?keyword=123123")
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.Len())
}
