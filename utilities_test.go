package qrcode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeqown/go-qrcode/v2/matrix"
)

func Test_min(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "min",
			args: args{
				x: 1,
				y: 2,
			},
			want: 1,
		},
		{
			name: "min",
			args: args{
				x: 2,
				y: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, min(tt.args.x, tt.args.y), "min(%v, %v)", tt.args.x, tt.args.y)
		})
	}
}

func Test_abs(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "abs",
			args: args{
				x: 1,
			},
			want: 1,
		},
		{
			name: "abs",
			args: args{
				x: -1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, abs(tt.args.x), "abs(%v)", tt.args.x)
		})
	}
}

func Test_samestate(t *testing.T) {
	type args struct {
		s1 matrix.State
		s2 matrix.State
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				s1: matrix.StateTrue,
				s2: matrix.StateTrue,
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				s1: matrix.StateFalse,
				s2: matrix.StateFalse,
			},
			want: true,
		},
		{
			name: "case 3",
			args: args{
				s1: matrix.StateTrue,
				s2: matrix.StateFalse,
			},
			want: false,
		},
		{
			name: "case 4",
			args: args{
				s1: matrix.StateFalse,
				s2: matrix.StateTrue,
			},
			want: false,
		},
		{
			name: "case 5",
			args: args{
				s1: matrix.StateFinder,
				s2: matrix.StateFinder,
			},
			want: true,
		},
		{
			name: "case 6",
			args: args{
				s1: matrix.StateFinder,
				s2: matrix.StateFalse,
			},
			want: false,
		},
		{
			name: "case 7",
			args: args{
				s1: matrix.StateTrue,
				s2: matrix.StateFinder,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, samestate(tt.args.s1, tt.args.s2), "samestate(%v, %v)", tt.args.s1, tt.args.s2)
		})
	}
}

func Benchmark_samestate(b *testing.B) {
	for i := 0; i < b.N; i++ {

		samestate(matrix.StateTrue, matrix.StateTrue)
		samestate(matrix.StateTrue, matrix.StateVersion)
	}
}

func Test_binaryToStateSlice(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []matrix.State
	}{
		{
			name: "case 1",
			args: args{
				"1010 0001 101a",
			},
			want: []matrix.State{
				// 1010
				matrix.StateTrue, matrix.StateFalse, matrix.StateTrue, matrix.StateFalse,
				// 0001
				matrix.StateFalse, matrix.StateFalse, matrix.StateFalse, matrix.StateTrue,
				// 101a
				matrix.StateTrue, matrix.StateFalse, matrix.StateTrue,
			},
		},
		{
			name: "case 2",
			args: args{
				"0000 11a1 11x2 x",
			},
			want: []matrix.State{
				// 0000
				matrix.StateFalse, matrix.StateFalse, matrix.StateFalse, matrix.StateFalse,
				// 11a1
				matrix.StateTrue, matrix.StateTrue, matrix.StateTrue,
				// 11x2
				matrix.StateTrue, matrix.StateTrue,
				// x
				// nothing
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, binaryToStateSlice(tt.args.s), "binaryToStateSlice(%v)", tt.args.s)
		})
	}
}
