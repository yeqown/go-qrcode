package qrcode

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewWith(t *testing.T) {
	qrc, err := NewWith("1234567", WithEncodingMode(EncModeNumeric),
		WithErrorCorrectionLevel(ErrorCorrectionLow))
	require.NoError(t, err)
	assert.NotNil(t, qrc)
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

// Test_NewWithReader exercises the
// io.Reader ability with all EncMode values
func Test_NewWithReader(t *testing.T) {

	okEncModes := []encMode{
		EncModeByte,
		EncModeNumeric,
	}

	panicEncModes := []encMode{
		EncModeAlphanumeric,
		EncModeJP,
	}

	// test EncMode that should be successful
	for _, em := range okEncModes {

		r := strings.NewReader("This is a wonderful QR code library")
		qrc, err := NewWithReader(r, WithEncodingMode(em),
			WithErrorCorrectionLevel(ErrorCorrectionLow))
		require.NoError(t, err)
		assert.NotNil(t, qrc)

		b := bytes.NewBuffer([]byte{8, 48, 37, 237, 187, 89, 0})
		qrc, err = NewWithReader(b, WithEncodingMode(em),
			WithErrorCorrectionLevel(ErrorCorrectionLow))
		require.NoError(t, err)
		assert.NotNil(t, qrc)
	}

	// test EncMode that are known to panic and its OK
	for _, em := range panicEncModes {

		assert.Panics(t, func() {
			r := strings.NewReader("This is a wonderful QR code library")
			_, _ = NewWithReader(r, WithEncodingMode(em),
				WithErrorCorrectionLevel(ErrorCorrectionLow))
		})

		assert.Panics(t, func() {
			b := bytes.NewBuffer([]byte{8, 48, 37, 237, 187, 89, 0})
			_, _ = NewWithReader(b, WithEncodingMode(em),
				WithErrorCorrectionLevel(ErrorCorrectionLow))
		})
	}
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
