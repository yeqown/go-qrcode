package qrcode

import (
	"testing"
)

// func init() {
// 	load(defaultVersionCfg)
// }

func TestEncodeNum(t *testing.T) {
	enc := encoder{
		ecLv:    ErrorCorrectionLow,
		mode:    EncModeNumeric,
		version: loadVersion(1, ErrorCorrectionLow),
	}

	b, err := enc.Encode([]byte("12312312"))
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

	b, err := enc.Encode([]byte("AKJA*:/"))
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

	b, err := enc.Encode([]byte("http://baidu.com?keyword=123123"))
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.Len())
}

func Test_analyzeNum(t *testing.T) {
	type args struct {
		byt byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{byt: '0'},
			want: true,
		},
		{
			name: "case 1",
			args: args{byt: 'a'},
			want: false,
		},
		{
			name: "case 2",
			args: args{byt: 'A'},
			want: false,
		},
		{
			name: "case 3",
			args: args{byt: '9'},
			want: true,
		},
		{
			name: "case 4",
			args: args{byt: '*'},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := analyzeNum(tt.args.byt); got != tt.want {
				t.Errorf("analyzeNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_analyzeAlphanum(t *testing.T) {
	type args struct {
		byt byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{byt: '0'},
			want: true,
		},
		{
			name: "case 1",
			args: args{byt: 'a'},
			want: false,
		},
		{
			name: "case 2",
			args: args{byt: 'A'},
			want: true,
		},
		{
			name: "case 3",
			args: args{byt: '9'},
			want: true,
		},
		{
			name: "case 4",
			args: args{byt: '*'},
			want: true,
		},
		{
			name: "case 5",
			args: args{byt: '?'},
			want: false,
		},
		{
			name: "case 6",
			args: args{byt: '&'},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := analyzeAlphaNum(tt.args.byt); got != tt.want {
				t.Errorf("analyzeAlphaNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_anlayzeMode(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want encMode
	}{
		{
			name: "case 0",
			args: args{raw: []byte("123120899231")},
			want: EncModeNumeric,
		},
		{
			name: "case 1",
			args: args{raw: []byte(":/1231H208*99231FBJO")},
			want: EncModeAlphanumeric,
		},
		{
			name: "case 2",
			args: args{raw: []byte("hahah1298312hG&^FBJO@jhgG*")},
			want: EncModeByte,
		},
		{
			name: "case 3",
			args: args{raw: []byte("JKAHDOIANKQOIHCMJKASJ")},
			want: EncModeAlphanumeric,
		},
		{
			name: "case 4",
			args: args{raw: []byte("https://baidu.com?keyword=_JSO==GA")},
			want: EncModeByte,
		},
		{
			name: "case 5",
			args: args{raw: []byte("这是汉字也应该是EncModeByte")},
			want: EncModeByte,
		},
		{
			name: "case 6 (swedish letter)",
			args: args{raw: []byte("Övrigt aksldjlk Övrigt should JP encMode?")},
			want: EncModeByte,
		},
		{
			name: "issue#28",
			args: args{raw: []byte("a")},
			want: EncModeByte,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := analyzeEncodeModeFromRaw(tt.args.raw); got != tt.want {
				t.Errorf("analyzeEncodeModeFromRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}
