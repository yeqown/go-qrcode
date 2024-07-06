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

func Test_analyzeJP(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{r: '0'},
			want: false,
		},
		{
			name: "case 1",
			args: args{r: 'a'},
			want: false,
		},
		{
			name: "case 2",
			args: args{r: 'A'},
			want: false,
		},
		{
			name: "case 3",
			args: args{r: '9'},
			want: false,
		},
		{
			name: "case 4",
			args: args{r: '*'},
			want: false,
		},
		{
			name: "case 5",
			args: args{r: '?'},
			want: false,
		},
		{
			name: "case 6",
			args: args{r: '&'},
			want: false,
		},
		{
			name: "case 7",
			args: args{r: 'Ö'},
			want: false,
		},
		{
			name: "case 8",
			args: args{r: 'に'},
			want: true,
		},
		{
			name: "case 9",
			args: args{r: '茗'},
			want: true,
		},
		{
			name: "case 10",
			args: args{r: '杆'},
			want: true,
		},
		{
			name: "case 11",
			args: args{r: '荷'},
			want: true,
		},
		{
			name: "case 12",
			args: args{r: '杠'},
			want: true,
		},
		{
			name: "case 13",
			args: args{r: '杙'},
			want: true,
		},
		{
			name: "case 14",
			args: args{r: '杣'},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := analyzeJP(tt.args.r); got != tt.want {
				t.Errorf("analyzeJP(%c=0x%x) = %v, want %v", tt.args.r, tt.args.r, got, tt.want)
			}
		})
	}
}

func Test_analyzeMode(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		want    encMode
		wantErr bool
	}{
		{
			name:    "case 0",
			args:    args{raw: "123120899231"},
			want:    EncModeNumeric,
			wantErr: false,
		},
		{
			name:    "case 1",
			args:    args{raw: ":/1231H208*99231FBJO"},
			want:    EncModeAlphanumeric,
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{raw: "hahah1298312hG&^FBJO@jhgG*"},
			want: EncModeByte,
		},
		{
			name:    "case 3",
			args:    args{raw: "JKAHDOIANKQOIHCMJKASJ"},
			want:    EncModeAlphanumeric,
			wantErr: false,
		},
		{
			name:    "case 4",
			args:    args{raw: "https://baidu.com?keyword=_JSO==GA"},
			want:    EncModeByte,
			wantErr: false,
		},
		{
			name:    "case 5",
			args:    args{raw: "茗荷"},
			want:    EncModeKanji,
			wantErr: false,
		},
		{
			name: "case 6 (swedish letter)",
			args: args{raw: "Övrigt aksldjlk Övrigt should JP encMode?"},
			want: EncModeByte,
		},
		{
			name:    "case 7 (japanese letter)",
			args:    args{raw: "朸 朷 杆 杞 杠 杙 杣"},
			want:    EncModeKanji,
			wantErr: false,
		},
		{
			name:    "issue#28",
			args:    args{raw: "a"},
			want:    EncModeByte,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := analyzeEncodeModeFromRaw(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("analyzeMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("analyzeMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
