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

// Test_NewWithConfig_UnmatchedEncodeMode tests that explicit encoding mode
// returns error when input contains characters that cannot be encoded.
func Test_NewWithConfig_UnmatchedEncodeMode(t *testing.T) {
	// Lowercase letters with Alphanumeric mode should return error
	_, err := NewWith("abcs", WithEncodingMode(EncModeAlphanumeric))
	assert.Error(t, err, "expected error when using lowercase letters with Alphanumeric mode")
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

// Test_NewWith_Kanji_EncMode tests Kanji mode encoding with explicit mode setting
func Test_NewWith_Kanji_EncMode(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		wantErr     bool
		expectedErr string
	}{
		// Valid Kanji input
		{
			name:    "single Kanji character",
			text:    "漢",
			wantErr: false,
		},
		{
			name:    "multiple Kanji characters",
			text:    "漢字",
			wantErr: false,
		},
		{
			name:    "Kanji sentence",
			text:    "日本語",
			wantErr: false,
		},
		{
			name:    "Kanji mixed characters",
			text:    "世界",
			wantErr: false,
		},
		// Invalid input for Kanji mode
		{
			name:        "ASCII characters",
			text:        "https://google.com",
			wantErr:     true,
			expectedErr: "cannot be encoded in kanji mode",
		},
		{
			name:        "numbers with Kanji mode",
			text:        "漢字123",
			wantErr:     true,
			expectedErr: "cannot be encoded in kanji mode",
		},
		{
			name:        "Hiragana with Kanji mode",
			text:        "こんにちは",
			wantErr:     true,
			expectedErr: "cannot be encoded in kanji mode",
		},
		{
			name:        "Katakana with Kanji mode",
			text:        "コンニチハ",
			wantErr:     true,
			expectedErr: "cannot be encoded in kanji mode",
		},
		{
			name:        "mixed Kanji and ASCII",
			text:        "漢字test",
			wantErr:     true,
			expectedErr: "cannot be encoded in kanji mode",
		},
		{
			name:        "CJK Extension A character",
			text:        "㐀",
			wantErr:     true,
			expectedErr: "cannot be encoded in kanji mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qrc, err := NewWith(tt.text,
				WithEncodingMode(EncModeKanji),
				WithErrorCorrectionLevel(ErrorCorrectionLow),
			)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, qrc)
			assert.Equal(t, EncModeKanji, qrc.encoder.mode)
			t.Logf("Kanji QR code for '%s': version=%d", tt.text, qrc.v.Ver)
		})
	}
}

// Test_NewWith_Kanji_AutoMode tests automatic Kanji mode detection
func Test_NewWith_Kanji_AutoMode(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected encMode
	}{
		{
			name:     "Kanji only - auto detect",
			text:     "漢字",
			expected: EncModeKanji,
		},
		{
			name:     "Kanji characters - auto detect",
			text:     "世界",
			expected: EncModeKanji,
		},
		{
			name:     "Mixed ASCII and Katakana - auto detect Byte mode",
			text:     "QRコード123",
			expected: EncModeByte,
		},
		{
			name:     "Pure Kanji text - auto detect",
			text:     "金木水火土日月星",
			expected: EncModeKanji,
		},
		{
			name:     "Long Kanji text - auto detect",
			text:     "東京京都大阪北海道沖縄鹿児島",
			expected: EncModeKanji,
		},
		{
			name:     "Hiragana only - auto detect Byte mode",
			text:     "これはひらがなです",
			expected: EncModeByte, // Hiragana is not supported in Kanji mode
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use EncModeAuto to let the library detect the mode automatically
			qrc, err := NewWith(tt.text,
				WithEncodingMode(EncModeAuto),
				WithErrorCorrectionLevel(ErrorCorrectionLow),
			)
			require.NoError(t, err)
			assert.NotNil(t, qrc)

			// Verify the detected mode matches expectation
			assert.Equal(t, tt.expected, qrc.encoder.mode,
				"Expected mode %v for text '%s', got %v", tt.expected, tt.text, qrc.encoder.mode)

			t.Logf("Auto-detected mode for '%s': %v, version=%d", tt.text, getEncModeName(qrc.encoder.mode), qrc.v.Ver)
		})
	}
}

// Test_NewWith_Kanji_Version10 tests Kanji encoding with specific version
func Test_NewWith_Kanji_Version10(t *testing.T) {
	qrc, err := NewWith("漢字文字試験",
		WithEncodingMode(EncModeKanji),
		WithVersion(10),
		WithErrorCorrectionLevel(ErrorCorrectionLow),
	)
	require.NoError(t, err)
	assert.NotNil(t, qrc)

	// Verify version is set correctly
	assert.Equal(t, 10, qrc.v.Ver)
	// Verify encoding mode is Kanji
	assert.Equal(t, EncModeKanji, qrc.encoder.mode)

	t.Logf("Kanji QR code with version 10: matrix dimension=%d", qrc.mat.Width())
}
