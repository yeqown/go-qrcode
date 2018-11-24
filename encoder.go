// Package qrcode ...
// encoder.go working for data encoding
package qrcode

import (
	"fmt"
	"log"

	"github.com/skip2/go-qrcode/bitset"
)

// EncMode ...
type EncMode uint

const (
	// EncModeNone mode ...
	EncModeNone EncMode = 1 << iota
	// EncModeNumeric mode ...
	EncModeNumeric
	// EncModeAlphanumeric mode ...
	EncModeAlphanumeric
	// EncModeByte mode ...
	EncModeByte
	// EncModeJP mode ...
	EncModeJP
)

var (
	paddingByte1 = bitset.NewFromBase2String("11101100")
	paddingByte2 = bitset.NewFromBase2String("00010001")
)

// GetEncModeName ...
func GetEncModeName(mode EncMode) string {
	switch mode {
	case EncModeNone:
		return "none"
	case EncModeNumeric:
		return "numeric"
	case EncModeAlphanumeric:
		return "alphanumeric"
	case EncModeByte:
		return "byte"
	case EncModeJP:
		return "japan"
	default:
		return "unknown"
	}
}

// getEncodeModeIndicator ...
func getEncodeModeIndicator(mode EncMode) *bitset.Bitset {
	switch mode {
	case EncModeNumeric:
		return bitset.New(false, false, false, true)
	case EncModeAlphanumeric:
		return bitset.New(false, false, true, false)
	case EncModeByte:
		return bitset.New(false, true, false, false)
	case EncModeJP:
		return bitset.New(true, false, false, false)
	default:
		panic("no indicator")
	}
}

// Encoder ... data to bit stream ...
type Encoder struct {
	// self init
	dst  *bitset.Bitset
	data []byte // raw input data

	// initial params
	mode EncMode // encode mode
	ecLv ECLevel // error correction level

	// self load
	version Version // QR verison ref
}

// Encode ...
// 1. encode raw data into bitset
// 2. append padding data
//
func (e *Encoder) Encode(byts []byte) (*bitset.Bitset, error) {
	e.dst = bitset.New()
	e.data = byts

	// appedn mode indicator symbol
	indicator := getEncodeModeIndicator(e.mode)
	e.dst.Append(indicator)
	// append chars length counter bits symbol
	e.dst.AppendUint32(uint32(len(byts)), e.charCountBits())

	// encode data with specified mode
	switch e.mode {
	case EncModeNumeric:
		e.encodeNumeric()
	case EncModeAlphanumeric:
		e.encodeAlphanumeric()
	case EncModeByte:
		e.encodeByte()
	case EncModeJP:
		panic("this has not been finished")
	}

	// fill and padding bits
	e.breakUpInto8bit()

	return e.dst, nil
}

// 0001b mode indicator
func (e *Encoder) encodeNumeric() {
	if e.dst == nil {
		log.Println("e.dst is nil")
		return
	}
	for i := 0; i < len(e.data); i += 3 {
		charsRemaining := len(e.data) - i

		var value uint32
		bitsUsed := 1

		for j := 0; j < charsRemaining && j < 3; j++ {
			value *= 10
			value += uint32(e.data[i+j] - 0x30)
			bitsUsed += 3
		}
		e.dst.AppendUint32(value, bitsUsed)
	}
}

// 0010b mode indicator
func (e *Encoder) encodeAlphanumeric() {
	if e.dst == nil {
		log.Println("e.dst is nil")
		return
	}
	for i := 0; i < len(e.data); i += 2 {
		charsRemaining := len(e.data) - i

		var value uint32
		for j := 0; j < charsRemaining && j < 2; j++ {
			value *= 45
			value += encodeAlphanumericCharacter(e.data[i+j])
		}

		bitsUsed := 6
		if charsRemaining > 1 {
			bitsUsed = 11
		}

		e.dst.AppendUint32(value, bitsUsed)
	}
}

// 0100b mode indicator
func (e *Encoder) encodeByte() {
	if e.dst == nil {
		log.Println("e.dst is nil")
		return
	}
	for _, b := range e.data {
		e.dst.AppendByte(b, 8)
	}
}

// Break Up into 8-bit Codewords and Add Pad Bytes if Necessary
func (e *Encoder) breakUpInto8bit() {
	// fill ending code (max 4bit)
	// depends on max capcity of current version and EC level
	maxCap := e.version.NumTotalCodewrods() * 8
	if less := maxCap - e.dst.Len(); less < 0 {
		panic("could not contain all char with wrong version cap")
	} else if less < 4 {
		e.dst.AppendNumBools(less, false)
	} else {
		e.dst.AppendNumBools(4, false)
	}

	// append `0` to be 8 times bits length
	if mod := e.dst.Len() % 8; mod != 0 {
		e.dst.AppendNumBools(8-mod, false)
	}

	// padding bytes
	// padding byte 11101100 00010001
	if n := maxCap - e.dst.Len(); n > 0 {
		debugLogf("maxCap: %d, len: %d, less: %d", maxCap, e.dst.Len(), n)
		for i := 1; i <= (n / 8); i++ {
			if i%2 == 1 {
				e.dst.Append(paddingByte1)
			} else {
				e.dst.Append(paddingByte2)
			}
		}
	}
}

// 字符计数指示符位长字典
var charCountMap = map[string]int{
	"9_numeric":       10,
	"9_alphanumeric":  9,
	"9_byte":          8,
	"9_japan":         8,
	"26_numeric":      12,
	"26_alphanumeric": 11,
	"26_byte":         16,
	"26_japan":        10,
	"40_numeric":      14,
	"40_alphanumeric": 13,
	"40_byte":         16,
	"40_japan":        12,
}

// charCountBits
func (e *Encoder) charCountBits() int {
	var lv int
	if v := e.version.Ver; v <= 9 {
		lv = 9
	} else if v <= 26 {
		lv = 26
	} else {
		lv = 40
	}
	pos := fmt.Sprintf("%d_%s", lv, GetEncModeName(e.mode))
	return charCountMap[pos]
}

// v must be a QR Code defined alphanumeric character: 0-9, A-Z, SP, $%*+-./ or
// :. The characters are mapped to values in the range 0-44 respectively.
func encodeAlphanumericCharacter(v byte) uint32 {
	c := uint32(v)

	switch {
	case c >= '0' && c <= '9':
		// 0-9 encoded as 0-9.
		return c - '0'
	case c >= 'A' && c <= 'Z':
		// A-Z encoded as 10-35.
		return c - 'A' + 10
	case c == ' ':
		return 36
	case c == '$':
		return 37
	case c == '%':
		return 38
	case c == '*':
		return 39
	case c == '+':
		return 40
	case c == '-':
		return 41
	case c == '.':
		return 42
	case c == '/':
		return 43
	case c == ':':
		return 44
	default:
		log.Panicf("encodeAlphanumericCharacter() with non alphanumeric char %c", v)
	}

	return 0
}

// 如果输入字符串只包含数字（0-9），请使用数字编码模式。
// 在数字编码模式不适用的情况下，如果可以在字符索引表的左列中找到输入字符串中的所有字符，请使用字符编码模式。注意：小写字母不能使用字符编码模式。
// 在字符编码模式不适用的情况下，如果字符可以在ISO-8859-1字符集中找到，则使用字节编码模式。
func chooseMode(raw []byte) EncMode {
	return EncModeByte
}
