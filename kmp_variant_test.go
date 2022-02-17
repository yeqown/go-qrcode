package qrcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_kmp(t *testing.T) {
	src := binaryToQRValueSlice("11001010010111001010010101011100")

	type args struct {
		src     []qrvalue
		pattern []qrvalue
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
	}{
		{
			name: "test1",
			args: args{
				src:     src,
				pattern: binaryToQRValueSlice("0101"),
			},
			wantCount: 6,
		},
		{
			name: "test1",
			args: args{
				src:     src,
				pattern: binaryToQRValueSlice("1001010"),
			},
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantCount, kmp(tt.args.src, tt.args.pattern, nil), "kmp(%v, %v)", tt.args.src, tt.args.pattern)
		})
	}
}
