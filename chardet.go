package qrcode

import (
	"errors"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	ErrNotSupportCharacter = errors.New("character set not supported, please check your input data")
)

// chardet.go refer to https://github.com/chardet/chardet to detect input string's
// character set, to see any unsupported character encountered in the input string.

// analyzeEncFunc returns true is current byte matched in current mode,
// otherwise means you should use a bigger character set to check.
type analyzeEncFunc func(rune) bool

// analyzeEncodeModeFromRaw try to detect letter set of input data,
// so that encoder can determine which mode should be use.
// reference: https://en.wikipedia.org/wiki/QR_code
//
// case1: only numbers, use EncModeNumeric.
// case2: could not use EncModeNumeric, but can find them all in character mapping, use EncModeAlphanumeric.
// case3: could not use EncModeAlphanumeric, but can find them all Shift JIS character set, use EncModeKanji.
// case4: could not use EncModeKanji, use EncModeByte.
func analyzeEncodeModeFromRaw(raw string) (encMode, error) {
	var (
		analyzeFn analyzeEncFunc
		mode      = EncModeNone
	)

	getNextAnalyzeFn := func() analyzeEncFunc {
		switch mode {
		case EncModeNumeric:
			return analyzeNum
		case EncModeAlphanumeric:
			return analyzeAlphaNum
		case EncModeKanji:
			return analyzeJP
		case EncModeByte:
			return analyzeByte
		default:
		}

		return nil
	}

	next := func() bool {
		// switch to next mode and get next analyze function. if no more analyze function, return true.
		mode <<= 1
		analyzeFn = getNextAnalyzeFn()
		return analyzeFn == nil
	}

	next()

	// Loop to check each character in raw data,
	// from low mode to higher while current mode could bear the input data.
	for _, r := range raw {
	reAnalyze:
		// issue#28 @borislavone reports this bug.
		// FIXED(@yeqown): next encMode analyzeVersionAuto func did not check the previous byte,
		// add goto statement to reanalyze previous byte which can't be analyzed in last encMode.
		if pass := analyzeFn(r); pass {
			continue
		}

		if nomore := next(); nomore {
			break
		}

		goto reAnalyze
	}

	if mode > EncModeByte {
		// If the mode overflow the EncModeKanji, means we can't encode the input data.
		return EncModeNone, ErrNotSupportCharacter
	}

	return mode, nil
}

// analyzeNum is r in num encMode
func analyzeNum(r rune) bool {
	return r >= '0' && r <= '9'
}

// analyzeAlphaNum is r in alpha number
func analyzeAlphaNum(r rune) bool {
	if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') {
		return true
	}
	switch r {
	case ' ', '$', '%', '*', '+', '-', '.', '/', ':':
		return true
	}
	return false
}

// analyzeByte always return true, since byte (utf8) mode can encode all characters.
func analyzeByte(r rune) bool {
	return true
}

// analyzeJP checks if a character can be encoded in QR Code Kanji mode.
// A character is valid for Kanji mode if:
// 1. It is in the CJK Unified Ideographs block (U+4E00-U+9FFF)
// 2. It can be converted to Shift JIS
// 3. The resulting Shift JIS value is in the valid QR Code ranges:
//   - 0x8140-0x9FFC (first range)
//   - 0xE040-0xEBBF (second range)
func analyzeJP(r rune) bool {
	// Check if the character is in the CJK Unified Ideographs block
	// This is a quick pre-check to avoid unnecessary conversion attempts
	// U+4E00-U+9FFF: CJK Unified Ideographs
	// U+3400-U+4DBF: CJK Unified Ideographs Extension A
	// U+F900-U+FAFF: CJK Compatibility Ideographs
	isCJK := (r >= 0x4E00 && r <= 0x9FFF) ||
		(r >= 0x3400 && r <= 0x4DBF) ||
		(r >= 0xF900 && r <= 0xFAFF)

	if !isCJK {
		return false
	}

	// Try to convert the character to Shift JIS
	// If conversion fails, it's not a valid Kanji character for QR Code
	enc := japanese.ShiftJIS.NewEncoder()
	s2, _, err := transform.String(enc, string(r))
	if err != nil || len(s2) != 2 {
		return false
	}

	// Check if the resulting Shift JIS value is in the valid QR Code Kanji ranges
	data := []byte(s2)
	hi := uint16(data[0])
	lo := uint16(data[1])
	code := hi<<8 | lo

	// QR Code Kanji mode supports Shift JIS ranges:
	// 0x8140-0x9FFC and 0xE040-0xEBBF
	return (code >= 0x8140 && code <= 0x9FFC) || (code >= 0xE040 && code <= 0xEBBF)
}
