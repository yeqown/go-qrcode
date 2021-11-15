package qrcode

import (
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
