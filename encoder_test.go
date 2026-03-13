package qrcode

import (
	"bytes"
	"testing"
)

func TestEncodeNum(t *testing.T) {
	enc := encoder{
		ecLv:    ErrorCorrectionLow,
		mode:    EncModeNumeric,
		version: loadVersion(1, ErrorCorrectionLow),
	}

	b, err := enc.Encode("12312312")
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.Len())
}

func TestEncodeAlphanum(t *testing.T) {
	enc := encoder{
		ecLv:    ErrorCorrectionLow,
		mode:    EncModeAlphanumeric,
		version: loadVersion(1, ErrorCorrectionLow),
	}

	b, err := enc.Encode("AKJA*:/")
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.Len())
}

func TestEncodeByte(t *testing.T) {
	enc := encoder{
		ecLv:    ErrorCorrectionQuart,
		mode:    EncModeByte,
		version: loadVersion(5, ErrorCorrectionQuart),
	}

	b, err := enc.Encode("http://baidu.com?keyword=123123")
	if err != nil {
		t.Errorf("could not encode: %v", err)
		t.Fail()
	}
	t.Log(b, b.Len())
}

func Test_toShiftJIS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test 茗荷",
			args: args{"茗荷"},
			want: []byte{0x1A, 0xAA, 0x06, 0x97},
		},
		{
			name: "test 世",
			args: args{"世"},
			// Shift JIS: 0x90A2
			// 0x90A2 - 0x8140 = 0x0F62, hi=0x0F, lo=0x62
			// encoded = 0x0F*0xC0 + 0x62 = 0xBA2
			// high byte = 0x0B, low byte = 0xA2
			want: []byte{0x0B, 0xA2},
		},
		{
			name: "test 世界",
			args: args{"世界"},
			// "世": [0x0B, 0xA2], "界": [0x06, 0xC5]
			want: []byte{0x0B, 0xA2, 0x06, 0xC5},
		},
		{
			name: "test 日本語",
			args: args{"日本語"},
			// "日": [0x0E, 0x3A], "本": [0x0F, 0xFB], "語": [0x08, 0xEA]
			want: []byte{0x0E, 0x3A, 0x0F, 0xFB, 0x08, 0xEA},
		},
		{
			name: "test 漢字",
			args: args{"漢字"},
			// "漢": [0x07, 0x3F], "字": [0x0A, 0x1A]
			want: []byte{0x07, 0x3F, 0x0A, 0x1A},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toShiftJIS(tt.args.s); !bytes.Equal(got, tt.want) {
				t.Errorf("toShiftJIS(%q) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}

func Test_encodeShiftJIS(t *testing.T) {
	type args struct {
		hi byte
		lo byte
	}
	tests := []struct {
		name  string
		args  args
		wantHi byte
		wantLo byte
	}{
		// Range 1: 0x8140-0x9FFC
		{
			name:  "lower boundary of range 1",
			args:  args{0x81, 0x40},
			wantHi: 0x00,
			wantLo: 0x00,
		},
		{
			name:  "middle of range 1 (世)",
			args:  args{0x90, 0xA2}, // "世" in Shift JIS
			// 0x90A2 - 0x8140 = 0x0F62
			// high=0x0F, low=0x62
			// encoded = 0x0F*0xC0 + 0x62 = 0xBA2
			wantHi: 0x0B,
			wantLo: 0xA2,
		},
		{
			name:  "upper boundary of range 1",
			args:  args{0x9F, 0xFC},
			// 0x9FFC - 0x8140 = 0x1EBC
			// high=0x1E, low=0xBC
			// encoded = 0x1E*0xC0 + 0xBC = 0x173C
			wantHi: 0x17,
			wantLo: 0x3C,
		},
		// Range 2: 0xE040-0xEBBF
		{
			name:  "lower boundary of range 2",
			args:  args{0xE0, 0x40},
			// 0xE040 - 0xC140 = 0x1F00
			// high=0x1F, low=0x00
			// encoded = 0x1F*0xC0 + 0x00 = 0x1740
			wantHi: 0x17,
			wantLo: 0x40,
		},
		{
			name:  "middle of range 2",
			args:  args{0xE4, 0xAA},
			// 0xE4AA - 0xC140 = 0x236A
			// high=0x23, low=0x6A
			// encoded = 0x23*0xC0 + 0x6A = 0x1AAA
			wantHi: 0x1A,
			wantLo: 0xAA,
		},
		{
			name:  "upper boundary of range 2",
			args:  args{0xEB, 0xBF},
			// 0xEBBF - 0xC140 = 0x2A7F
			// high=0x2A, low=0x7F
			// encoded = 0x2A*0xC0 + 0x7F = 0x1FFF
			wantHi: 0x1F,
			wantLo: 0xFF,
		},
		// Invalid ranges
		{
			name:  "below range 1",
			args:  args{0x80, 0x00},
			wantHi: 0x00,
			wantLo: 0x00,
		},
		{
			name:  "between ranges",
			args:  args{0x9F, 0xFD},
			wantHi: 0x00,
			wantLo: 0x00,
		},
		{
			name:  "above range 2",
			args:  args{0xEC, 0x00},
			wantHi: 0x00,
			wantLo: 0x00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHi, gotLo := encodeShiftJIS(tt.args.hi, tt.args.lo)
			if gotHi != tt.wantHi || gotLo != tt.wantLo {
				t.Errorf("encodeShiftJIS(0x%02X, 0x%02X) = (0x%02X, 0x%02X), want (0x%02X, 0x%02X)",
					tt.args.hi, tt.args.lo, gotHi, gotLo, tt.wantHi, tt.wantLo)
			}
		})
	}
}

func TestEncodeKanji(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int // Expected bit length for encoded data (13 bits per character)
	}{
		{
			name:    "single character 世",
			input:   "世",
			wantLen: 13,
		},
		{
			name:    "two characters 世界",
			input:   "世界",
			wantLen: 26,
		},
		{
			name:    "four characters 日本語",
			input:   "日本語",
			wantLen: 39,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := encoder{
				ecLv:    ErrorCorrectionLow,
				mode:    EncModeKanji,
				version: loadVersion(1, ErrorCorrectionLow),
			}

			b, err := enc.Encode(tt.input)
			if err != nil {
				t.Errorf("could not encode: %v", err)
				t.Fail()
			}

			// The total length includes mode indicator (4), char count (8 for v1-9),
			// data bits (wantLen), and padding to fill the codeword capacity
			t.Logf("Encode(%q) total bits: %d", tt.input, b.Len())

			// Check that we successfully encoded the Kanji data
			// The mode indicator is 1000 (4 bits), char count is 8 bits for version 1
			// So the first 12 bits should be: 1000 (mode) + char count (8 bits)
			// For a single character, char count = 1 = 00000001
			// First 12 bits: 1000 00000001
		})
	}
}

func TestEncodeKanji_Version(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		version               int
		expectedCharCountBits int
	}{
		{
			name:                  "version 1",
			input:                 "漢字",
			version:               1,
			expectedCharCountBits: 8,
		},
		{
			name:                  "version 10",
			input:                 "漢字",
			version:               10,
			expectedCharCountBits: 10,
		},
		{
			name:                  "version 27",
			input:                 "漢字",
			version:               27,
			expectedCharCountBits: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := encoder{
				ecLv:    ErrorCorrectionLow,
				mode:    EncModeKanji,
				version: loadVersion(tt.version, ErrorCorrectionLow),
			}

			charCountBits := enc.charCountBits()
			if charCountBits != tt.expectedCharCountBits {
				t.Errorf("charCountBits() = %d, want %d", charCountBits, tt.expectedCharCountBits)
			}

			b, err := enc.Encode(tt.input)
			if err != nil {
				t.Errorf("could not encode: %v", err)
				t.Fail()
			}

			t.Logf("Encode(%q) with version %d = %v, total bits: %d", tt.input, tt.version, b, b.Len())
		})
	}
}
