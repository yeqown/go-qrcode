package qrcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_qrtype(t *testing.T) {
	assert.Equal(t, uint8(0b00000010), uint8(QRType_INIT))     // 1 << 1
	assert.Equal(t, uint8(0b00000100), uint8(QRType_DATA))     // 2 << 1
	assert.Equal(t, uint8(0b00000110), uint8(QRType_VERSION))  // 3 << 1
	assert.Equal(t, uint8(0b00001000), uint8(QRType_FORMAT))   // 4 << 1
	assert.Equal(t, uint8(0b00001010), uint8(QRType_FINDER))   // 5 << 1
	assert.Equal(t, uint8(0b00001100), uint8(QRType_DARK))     // 6 << 1
	assert.Equal(t, uint8(0b00001110), uint8(QRType_SPLITTER)) // 7 << 1
	assert.Equal(t, uint8(0b00010000), uint8(QRType_TIMING))   // 8 << 1

}

func Test_qrvalue(t *testing.T) {
	// QRValue_INIT_V0
	assert.Equal(t, QRType_INIT, QRValue_INIT_V0.qrtype())
	assert.False(t, QRValue_INIT_V0.qrbool())

	// QRValue_DATA_V1
	assert.Equal(t, QRType_DATA, QRValue_DATA_V1.qrtype())
	assert.True(t, QRValue_DATA_V1.qrbool())

	// QRValue_DATA_V0
	assert.Equal(t, QRType_DATA, QRValue_DATA_V0.qrtype())
	assert.False(t, QRValue_DATA_V0.qrbool())

	// QRValue_VERSION_V0
	assert.Equal(t, QRType_VERSION, QRValue_VERSION_V0.qrtype())
	assert.False(t, QRValue_VERSION_V0.qrbool())

	// QRValue_VERSION_V1
	assert.Equal(t, QRType_VERSION, QRValue_VERSION_V1.qrtype())
	assert.True(t, QRValue_VERSION_V1.qrbool())

	// QRValue_FORMAT_V0
	assert.Equal(t, QRType_FORMAT, QRValue_FORMAT_V0.qrtype())
	assert.False(t, QRValue_FORMAT_V0.qrbool())

	// QRValue_FORMAT_V1
	assert.Equal(t, QRType_FORMAT, QRValue_FORMAT_V1.qrtype())
	assert.True(t, QRValue_FORMAT_V1.qrbool())

	// QRValue_FINDER_V0
	assert.Equal(t, QRType_FINDER, QRValue_FINDER_V0.qrtype())
	assert.False(t, QRValue_FINDER_V0.qrbool())

	// QRValue_FINDER_V1
	assert.Equal(t, QRType_FINDER, QRValue_FINDER_V1.qrtype())
	assert.True(t, QRValue_FINDER_V1.qrbool())

	// QRValue_DARK_V0
	assert.Equal(t, QRType_DARK, QRValue_DARK_V0.qrtype())
	assert.False(t, QRValue_DARK_V0.qrbool())

	// QRValue_DARK_V1
	assert.Equal(t, QRType_DARK, QRValue_DARK_V1.qrtype())
	assert.True(t, QRValue_DARK_V1.qrbool())

	// QRValue_SPLITTER_V0
	assert.Equal(t, QRType_SPLITTER, QRValue_SPLITTER_V0.qrtype())
	assert.False(t, QRValue_SPLITTER_V0.qrbool())

	// QRValue_SPLITTER_V1
	assert.Equal(t, QRType_SPLITTER, QRValue_SPLITTER_V1.qrtype())
	assert.True(t, QRValue_SPLITTER_V1.qrbool())

	// QRValue_TIMING_V0
	assert.Equal(t, QRType_TIMING, QRValue_TIMING_V0.qrtype())
	assert.False(t, QRValue_TIMING_V0.qrbool())

	// QRValue_TIMING_V1
	assert.Equal(t, QRType_TIMING, QRValue_TIMING_V1.qrtype())
	assert.True(t, QRValue_TIMING_V1.qrbool())
}

func Test_qrvalue_xor(t *testing.T) {
	type args struct {
		s1 qrvalue
		s2 qrvalue
	}
	tests := []struct {
		name string
		args args
		want qrvalue
	}{
		{
			name: "case1",
			args: args{
				s1: QRValue_DATA_V0,
				s2: QRValue_DATA_V0,
			},
			want: QRValue_DATA_V0,
		},
		{
			name: "case1",
			args: args{
				s1: QRValue_DATA_V1,
				s2: QRValue_DATA_V0,
			},
			want: QRValue_DATA_V1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.s1.xor(tt.args.s2); got != tt.want {
				t.Errorf("XOR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_qrtype_String(t *testing.T) {
	tests := []struct {
		name string
		t    qrtype
		want string
	}{
		{
			name: "case1",
			t:    QRType_INIT,
			want: "I",
		},
		{
			name: "case2",
			t:    QRType_DATA,
			want: "d",
		},
		{
			name: "case3",
			t:    QRType_VERSION,
			want: "V",
		},
		{
			name: "case4",
			t:    QRType_FORMAT,
			want: "f",
		},
		{
			name: "case5",
			t:    QRType_FINDER,
			want: "F",
		},
		{
			name: "case6",
			t:    QRType_DARK,
			want: "D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("qrtype.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_qrvalue_String(t *testing.T) {
	tests := []struct {
		name string
		s    qrvalue
		want string
	}{
		{
			name: "data 0",
			s:    QRValue_DATA_V0,
			want: "d0",
		},
		{
			name: "data 1",
			s:    QRValue_DATA_V1,
			want: "d1",
		},
		{
			name: "version 0",
			s:    QRValue_VERSION_V0,
			want: "V0",
		},
		{
			name: "version 1",
			s:    QRValue_VERSION_V1,
			want: "V1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("qrvalue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
