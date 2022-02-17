package qrcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		s1 qrvalue
		s2 qrvalue
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				s1: QRValue_DATA_V1,
				s2: QRValue_DATA_V1,
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				s1: QRValue_DATA_V0,
				s2: QRValue_DATA_V0,
			},
			want: true,
		},
		{
			name: "case 3",
			args: args{
				s1: QRValue_DATA_V1,
				s2: QRValue_DATA_V0,
			},
			want: false,
		},
		{
			name: "case 4",
			args: args{
				s1: QRValue_DATA_V0,
				s2: QRValue_DATA_V1,
			},
			want: false,
		},
		{
			name: "case 5",
			args: args{
				s1: QRValue_FINDER_V1,
				s2: QRValue_FINDER_V1,
			},
			want: true,
		},
		{
			name: "case 6",
			args: args{
				s1: QRValue_FINDER_V1,
				s2: QRValue_DATA_V0,
			},
			want: false,
		},
		{
			name: "case 7",
			args: args{
				s1: QRValue_DATA_V1,
				s2: QRValue_FINDER_V1,
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
		samestate(QRValue_DATA_V1, QRValue_DATA_V1)
		samestate(QRValue_DATA_V1, QRValue_VERSION_V1)
	}
}

func Test_binaryToStateSlice(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []qrvalue
	}{
		{
			name: "case 1",
			args: args{
				"1010 0001 101a",
			},
			want: []qrvalue{
				// 1010
				QRValue_DATA_V1, QRValue_DATA_V0, QRValue_DATA_V1, QRValue_DATA_V0,
				// 0001
				QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V1,
				// 101a
				QRValue_DATA_V1, QRValue_DATA_V0, QRValue_DATA_V1,
			},
		},
		{
			name: "case 2",
			args: args{
				"0000 11a1 11x2 x",
			},
			want: []qrvalue{
				// 0000
				QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0,
				// 11a1
				QRValue_DATA_V1, QRValue_DATA_V1, QRValue_DATA_V1,
				// 11x2
				QRValue_DATA_V1, QRValue_DATA_V1,
				// x
				// nothing
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, binaryToQRValueSlice(tt.args.s), "binaryToQRValueSlice(%v)", tt.args.s)
		})
	}
}
