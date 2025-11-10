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

// Test_NewWith_MinimumVersion tests the WithMinimumVersion option
func Test_NewWith_MinimumVersion(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		minVersion     int
		expectedMinVer int
	}{
		{
			name:           "short text with minimum version 9",
			text:           "123",
			minVersion:     9,
			expectedMinVer: 9,
		},
		{
			name:           "short text with minimum version 15",
			text:           "HELLO",
			minVersion:     15,
			expectedMinVer: 15,
		},
		{
			name:           "medium text naturally version 5 with minimum version 9",
			text:           strings.Repeat("A", 100),
			minVersion:     9,
			expectedMinVer: 9,
		},
		{
			name:           "long text naturally version 20 with minimum version 9",
			text:           strings.Repeat("A", 500),
			minVersion:     9,
			expectedMinVer: 9, // should be at least 9 (actual may be higher)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qrc, err := NewWith(tt.text,
				WithMinimumVersion(tt.minVersion),
				WithErrorCorrectionLevel(ErrorCorrectionQuart),
			)
			require.NoError(t, err)
			assert.NotNil(t, qrc)
			assert.GreaterOrEqual(t, qrc.v.Ver, tt.expectedMinVer,
				"QR code version should be at least %d, got %d", tt.expectedMinVer, qrc.v.Ver)
		})
	}
}

// Test_NewWith_MinimumVersion_InvalidValues tests that invalid minimum versions are ignored
func Test_NewWith_MinimumVersion_InvalidValues(t *testing.T) {
	tests := []struct {
		name       string
		minVersion int
		text       string
	}{
		{
			name:       "minimum version 0 should be ignored",
			minVersion: 0,
			text:       "123",
		},
		{
			name:       "minimum version 41 should be ignored",
			minVersion: 41,
			text:       "123",
		},
		{
			name:       "minimum version -1 should be ignored",
			minVersion: -1,
			text:       "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qrc, err := NewWith(tt.text,
				WithMinimumVersion(tt.minVersion),
				WithErrorCorrectionLevel(ErrorCorrectionLow),
			)
			require.NoError(t, err)
			assert.NotNil(t, qrc)
			// Should still create a valid QR code even if minimum version was invalid
			assert.GreaterOrEqual(t, qrc.v.Ver, 1)
			assert.LessOrEqual(t, qrc.v.Ver, 40)
		})
	}
}

// Test_NewWith_MinimumVersion_WithExplicitVersion tests interaction between Version and MinimumVersion
func Test_NewWith_MinimumVersion_WithExplicitVersion(t *testing.T) {
	// When both WithVersion and WithMinimumVersion are set, WithVersion takes precedence
	qrc, err := NewWith("123",
		WithVersion(10),
		WithMinimumVersion(15),
		WithErrorCorrectionLevel(ErrorCorrectionLow),
	)
	require.NoError(t, err)
	assert.NotNil(t, qrc)
	// WithVersion takes precedence, so version should be 10
	assert.Equal(t, 10, qrc.v.Ver)
}
