package qrcode

import (
	"log"
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
// case2: could not use EncModeNumeric, but you can find all of them in character mapping, use EncModeAlphanumeric.
// case3: could not use EncModeAlphanumeric, but you can find all of them in ISO-8859-1 character set, use EncModeByte.
// case4: could not use EncModeByte, use EncModeJP, no more choice.
func analyzeEncodeModeFromRaw(raw string) encMode {
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
		case EncModeByte:
			return analyzeByte
		case EncModeJP:
			return analyzeJP
		default:
		}

		return analyzeDefault
	}

	next := func() {
		// switch to next mode and get next analyze function.
		mode <<= 1
		analyzeFn = getNextAnalyzeFn()
	}

	next()

	// Loop to check each character in raw data,
	// from low mode to higher while current mode could bear the input data.
	for _, byt := range raw {
	reAnalyze:
		// issue#28 @borislavone reports this bug.
		// FIXED(@yeqown): next encMode analyzeVersionAuto func did not check the previous byte,
		// add goto statement to reanalyze previous byte which can't be analyzed in last encMode.
		if !analyzeFn(byt) {
			next()
			goto reAnalyze
		}
	}

	if mode > EncModeJP {
		// If the mode overflow the EncModeJP, means we can't encode the input data.
		log.Panicf("could not encode the input data: %s", raw)
	}

	return mode
}

// analyzeNum is byt in num encMode
func analyzeNum(r rune) bool {
	return r >= '0' && r <= '9'
}

// analyzeAlphaNum is byt in alpha number
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

// analyzeByte contains ISO-8859-1 character set
func analyzeByte(r rune) bool {
	// ISO-8859-1 character set, if r > \u00ff, means it's not in ISO-8859-1.
	return r <= '\u00ff'
}

// analyzeJP contains Kanji character set
// http://www.rikai.com/library/kanjitables/kanji_codes.sjis.shtml
func analyzeJP(r rune) bool {
	// Kanji character set
	if r > 0x8140 && r < 0x9FFC {
		return true
	}
	if r > 0xE040 && r < 0xEBBF {
		return true
	}

	return false
}

func analyzeDefault(r rune) bool {
	return false
}
