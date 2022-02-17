package qrcode

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewWith(t *testing.T) {
	qrc, err := NewWith("1234567",
		WithEncodingMode(EncModeNumeric),
		WithErrorCorrectionLevel(ErrorCorrectionLow),
		WithVersion(7),
	)
	require.NoError(t, err)
	assert.NotNil(t, qrc)

	SetDebugMode()
	_ = debugDraw("./testdata/nw.jpeg", *qrc.mat.Copy())

	qrc.mat.print()
}

// Test_NewWithConfig_UnmatchedEncodeMode NewWith will panic while encMode is
// not matched to Config.EncMode, for example:
// cfg.EncMode is EncModeAlphanumeric but source text is bytes encoding.
func Test_NewWithConfig_UnmatchedEncodeMode(t *testing.T) {
	assert.Panics(t, func() {
		_, err := NewWith("abcs", WithEncodingMode(EncModeAlphanumeric))
		if err != nil {
			t.Errorf("could not generate QRCode: %v", err)
			t.Fail()
		}
	})
}

func Benchmark_NewQRCode_1KB(b *testing.B) {
	text := strings.Repeat("abcdefghij", 100)

	for i := 0; i < b.N; i++ {
		_, err := New(text)
		if err != nil {
			b.Errorf("could not generate QRCode: %v", err)
			b.Fail()
		}
	}
}
