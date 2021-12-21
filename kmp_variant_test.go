package qrcode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeqown/go-qrcode/v2/matrix"
)

func Test_kmp(t *testing.T) {
	src := binaryToStateSlice("11001010010111001010010101011100")

	type args struct {
		src     []matrix.State
		pattern []matrix.State
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
				pattern: binaryToStateSlice("0101"),
			},
			wantCount: 6,
		},
		{
			name: "test1",
			args: args{
				src:     src,
				pattern: binaryToStateSlice("1001010"),
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
