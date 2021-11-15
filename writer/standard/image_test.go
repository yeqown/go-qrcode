package standard

import (
	"image/color"
	"os"
	"reflect"
	"testing"

	"github.com/yeqown/go-qrcode/v2/matrix"

	"github.com/stretchr/testify/require"
)

func Test_image_draw(t *testing.T) {
	m := matrix.New(20, 20)
	// set all 3rd column as black else be white
	for x := 0; x < m.Width(); x++ {
		_ = m.Set(x, 3, matrix.StateTrue)
	}

	fd, err := os.Create("./testdata/default.jpeg")
	require.NoError(t, err)
	err = drawTo(fd, *m, nil)
	require.NoError(t, err)
}

func Test_stateRGBA(t *testing.T) {
	type args struct {
		v matrix.State
	}
	tests := []struct {
		name string
		args args
		want color.Color
	}{
		{
			name: "case 1",
			args: args{v: matrix.StateFalse},
			want: _stateToRGBA[matrix.StateFalse],
		},
		{
			name: "case 2",
			args: args{v: matrix.StateInit},
			want: _stateToRGBA[matrix.StateInit],
		},
		{
			name: "case 3",
			args: args{v: matrix.StateTrue},
			want: _stateToRGBA[matrix.StateTrue],
		},
		{
			name: "case 4",
			args: args{v: matrix.StateFormat},
			want: _defaultStateColor,
		},
		{
			name: "case 5",
			args: args{v: matrix.StateVersion},
			want: _defaultStateColor,
		},
		{
			name: "case 6",
			args: args{v: matrix.State(0x6767)},
			want: _defaultStateColor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := new(outputImageOptions).stateRGBA(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stateRGBA() = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_hexToRGBA(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want color.RGBA
	}{
		{
			name: "case 1",
			args: args{s: "#112233"},
			want: color.RGBA{R: 17, G: 34, B: 51, A: 255},
		},
		{
			name: "case 2",
			args: args{s: "#112"},
			want: color.RGBA{R: 17, G: 17, B: 34, A: 255},
		},
		//{
		//	name: "case 3",
		//	args: args{s: "#1122331"},
		//	want: color.RGBA{},
		//}, // panic
		//{
		//	name: "case 4",
		//	args: args{s: "#11"},
		//	want: color.RGBA{},
		//}, // panic
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hexToRGBA(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hexToRGBA() = %v, want %v", got, tt.want)
			}
		})
	}
}
