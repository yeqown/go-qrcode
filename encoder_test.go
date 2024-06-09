package qrcode

import (
	"bytes"
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

func Test_toShiftJIS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test 1",
			args: args{"茗荷"},
			want: []byte{0x1A, 0xAA, 0x06, 0x97},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toShiftJIS(tt.args.s); !bytes.Equal(got, tt.want) {
				t.Errorf("toShiftJIS() = %v, want %v", got, tt.want)
			}
		})
	}
}
