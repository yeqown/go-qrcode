package qrcode

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
	analyzeFnMapping := map[encMode]analyzeEncFunc{
		EncModeNumeric:      analyzeNum,
		EncModeAlphanumeric: analyzeAlphaNum,
		EncModeByte:         analyzeByte,
		EncModeJP:           nil,
	}

	var (
		analyzeFn analyzeEncFunc
		mode      = EncModeNumeric
	)

	// loop to check each character in raw data,
	// from low mode to higher while current mode could bearing the input data.
	for _, byt := range raw {
	reAnalyze:
		if analyzeFn = analyzeFnMapping[mode]; analyzeFn == nil {
			break
		}

		// issue#28 @borislavone reports this bug.
		// FIXED(@yeqown): next encMode analyzeVersionAuto func did not check the previous byte,
		// add goto statement to reanalyze previous byte which can't be analyzed in last encMode.
		if !analyzeFn(byt) {
			mode <<= 1
			goto reAnalyze
		}
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
	if r > '\u00ff' {
		return false
	}

	return true
}
