package qrcode

import (
	"testing"
)

func Test_analyzeNum(t *testing.T) {
	type args struct {
		byt rune
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
		byt rune
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

func Test_analyzeByte(t *testing.T) {
	type args struct {
		byt rune
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
			want: true,
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
			want: true,
		},
		{
			name: "case 6",
			args: args{byt: '&'},
			want: true,
		},
		{
			name: "case 7",
			args: args{byt: 'Ö'},
			want: true,
		},
		{
			name: "case 8",
			args: args{byt: 'に'},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := analyzeByte(tt.args.byt); got != tt.want {
				t.Errorf("analyzeByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_analyzeMode(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name string
		args args
		want encMode
	}{
		{
			name: "case 0",
			args: args{raw: "123120899231"},
			want: EncModeNumeric,
		},
		{
			name: "case 1",
			args: args{raw: ":/1231H208*99231FBJO"},
			want: EncModeAlphanumeric,
		},
		{
			name: "case 2",
			args: args{raw: "hahah1298312hG&^FBJO@jhgG*"},
			want: EncModeByte,
		},
		{
			name: "case 3",
			args: args{raw: "JKAHDOIANKQOIHCMJKASJ"},
			want: EncModeAlphanumeric,
		},
		{
			name: "case 4",
			args: args{raw: "https://baidu.com?keyword=_JSO==GA"},
			want: EncModeByte,
		},
		{
			name: "case 5",
			args: args{raw: "这是汉字也应该是EncModeByte"},
			want: EncModeJP,
		},
		{
			name: "case 6 (swedish letter)",
			args: args{raw: "Övrigt aksldjlk Övrigt should JP encMode?"},
			want: EncModeByte,
		},
		{
			name: "case 7 (japanese letter)",
			args: args{raw: "にほんごのテスト"},
			want: EncModeJP,
		},
		{
			name: "issue#28",
			args: args{raw: "a"},
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
