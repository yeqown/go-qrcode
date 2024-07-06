// Package qrcode ...
// encoder.go working for data encoding
package qrcode

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yeqown/reedsolomon/binary"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// encMode indicates the encoding mode of the data to be encoded.
// The encoding mode is used to determine how the data should be encoded
// into bits for the QR code. This repository supports the following encoding
// modes:
// - EncModeNone: no encoding
// - EncModeNumeric: numeric encoding
// - EncModeAlphanumeric: alphanumeric encoding
// - EncModeByte: byte encoding
// - EncModeKanji: japanese encoding
//
// The encoding mode is determined by the data to be encoded. For example, if
// the data to be encoded is all numeric, the encoding mode will be EncModeNumeric.
// If the data to be encoded is alphanumeric, the encoding mode will be EncModeAlphanumeric.
// You can also specify the encoding mode automatically by using EncModeAuto, which
// will automatically determine the encoding mode based on the data to be encoded.
type encMode uint

const (
	// EncModeAuto will trigger a detection of the letter set from the input data.
	EncModeAuto = 0
	// EncModeNone mode ...
	EncModeNone encMode = 1 << iota
	// EncModeNumeric mode ...
	EncModeNumeric
	// EncModeAlphanumeric mode ...
	EncModeAlphanumeric
	// EncModeByte mode ...
	EncModeByte
	// EncModeJP mode ...
	// @Deprecated use EncModeKanji instead
	EncModeJP
	EncModeKanji = EncModeJP
)

var (
	paddingByte1, _ = binary.NewFromBinaryString("11101100")
	paddingByte2, _ = binary.NewFromBinaryString("00010001")
)

// getEncModeName ...
func getEncModeName(mode encMode) string {
	switch mode {
	case EncModeNone:
		return "none"
	case EncModeNumeric:
		return "numeric"
	case EncModeAlphanumeric:
		return "alphanumeric"
	case EncModeByte:
		return "byte"
	case EncModeKanji:
		return "kanji"
	default:
		return "unknown(" + strconv.Itoa(int(mode)) + ")"
	}
}

// getEncodeModeIndicator ...
func getEncodeModeIndicator(mode encMode) *binary.Binary {
	switch mode {
	case EncModeNumeric:
		return binary.New(false, false, false, true)
	case EncModeAlphanumeric:
		return binary.New(false, false, true, false)
	case EncModeByte:
		return binary.New(false, true, false, false)
	case EncModeKanji:
		return binary.New(true, false, false, false)
	default:
		panic("no indicator")
	}
}

// encoder ... data to bit stream ...
type encoder struct {
	// self init
	dst *binary.Binary

	// initial params
	mode encMode // encode mode
	ecLv ecLevel // error correction level

	// self load
	version version // QR version ref
}

func newEncoder(m encMode, ec ecLevel, v version) *encoder {
	switch m {
	case EncModeNumeric, EncModeAlphanumeric, EncModeByte, EncModeKanji:
	default:
		panic("unsupported data encoding mode in newEncoder()")
	}

	return &encoder{
		dst:     nil,
		mode:    m,
		ecLv:    ec,
		version: v,
	}
}

// Encode ...
// 1. encode raw data into bitset
// 2. append _defaultPadding data
func (e *encoder) Encode(raw string) (*binary.Binary, error) {
	e.dst = binary.New()

	var data []byte
	switch e.mode {
	case EncModeNumeric, EncModeAlphanumeric, EncModeByte:
		data = []byte(raw)
	case EncModeKanji:
		data = toShiftJIS(raw)
	default:
		log.Printf("unsupported encoding mode: %s", getEncModeName(e.mode))
	}

	// append mode indicator symbol
	indicator := getEncodeModeIndicator(e.mode)
	e.dst.Append(indicator)
	// append chars length counter bits symbol
	e.dst.AppendUint32(uint32(len(data)), e.charCountBits())

	// encode data with specified mode
	switch e.mode {
	case EncModeNumeric:
		e.encodeNumeric(data)
	case EncModeAlphanumeric:
		e.encodeAlphanumeric(data)
	case EncModeByte:
		e.encodeByte(data)
	case EncModeKanji:
		e.encodeKanji(data)
	default:
		log.Printf("unsupported encoding mode: %s", getEncModeName(e.mode))
	}

	// fill and _defaultPadding bits
	e.breakUpInto8bit()

	return e.dst, nil
}

// 0001b mode indicator
func (e *encoder) encodeNumeric(data []byte) {
	if e.dst == nil {
		log.Println("e.dst is nil")
		return
	}
	for i := 0; i < len(data); i += 3 {
		charsRemaining := len(data) - i

		var value uint32
		bitsUsed := 1

		for j := 0; j < charsRemaining && j < 3; j++ {
			value *= 10
			value += uint32(data[i+j] - 0x30)
			bitsUsed += 3
		}
		e.dst.AppendUint32(value, bitsUsed)
	}
}

// 0010b mode indicator
func (e *encoder) encodeAlphanumeric(data []byte) {
	if e.dst == nil {
		log.Println("e.dst is nil")
		return
	}
	for i := 0; i < len(data); i += 2 {
		charsRemaining := len(data) - i

		var value uint32
		for j := 0; j < charsRemaining && j < 2; j++ {
			value *= 45
			value += encodeAlphanumericCharacter(data[i+j])
		}

		bitsUsed := 6
		if charsRemaining > 1 {
			bitsUsed = 11
		}

		e.dst.AppendUint32(value, bitsUsed)
	}
}

// 0100b mode indicator
func (e *encoder) encodeByte(data []byte) {
	if e.dst == nil {
		log.Println("e.dst is nil")
		return
	}
	for _, b := range data {
		_ = e.dst.AppendByte(b, 8)
	}
}

// toShiftJIS
// https://www.thonky.com/qr-code-tutorial/kanji-mode-encoding
func toShiftJIS(raw string) []byte {
	// FIXME: some character encoded into Shift JIS but not in the range of 0x8140-0x9FFC and 0xE040-0xEBBF.
	enc := japanese.ShiftJIS.NewEncoder()
	s2, _, err := transform.String(enc, raw)
	if err != nil {
		log.Printf("could not encode string to Shift JIS: %v", err)
		return []byte{}
	}

	data := []byte(s2)
	if len(data)%2 != 0 {
		// BUG: encode bytes with Shift JIS must be times of 2, cause panic here
		log.Panicf("shift JIS encoded []byte must be times of 2, but got %d", len(data))
	}

	for i := 0; i < len(data); i += 2 {
		data[i], data[i+1] = encodeShiftJIS(data[i], data[i+1])
	}

	return data
}

func encodeShiftJIS(hi byte, lo byte) (byte, byte) {
	r := uint16(hi)<<8 | uint16(lo)

	// fmt.Printf("before: r=%x\n", r)
	if r > 0x8140 && r < 0x9FFC {
		r -= 0x8140
	} else if r > 0xE040 && r < 0xEBBF {
		r -= 0xC140
	} else {
		// Not a Shift JIS character out of range 0x8140-0x9FFC and 0xE040-0xEBBF
		log.Printf("'%c'(0x%x) not a Shift JIS character out of range 0x8140-0x9FFC and 0xE040-0xEBBF", r, r)
		return 0, 0
	}

	fmt.Printf("middle: r=%x\n", r)
	hi = uint8(r >> 8)
	lo = uint8(r & 0xFF)

	// fmt.Printf("middle: high=%x, low=%x\n", hi, lo)

	r = uint16(hi)*uint16(0xC0) + uint16(lo)
	// fmt.Printf("after: r=%x\n", r)

	return byte(r >> 8), byte(r & 0xFF)
}

// encodeKanji
func (e *encoder) encodeKanji(data []byte) {
	// data must be times of 2, since toShiftJIS encode 1 char to 2 bytes
	if len(data)%2 != 0 {
		log.Println("data must be times of 2")
	}

	for i := 0; i < len(data); i += 2 {
		// 2 bytes to 1 kanji
		// 2 bytes to 13 bits
		_ = e.dst.AppendByte(data[i]<<3, 5)
		_ = e.dst.AppendByte(data[i+1], 8)
	}
}

// Break Up into 8-bit Codewords and Add Pad Bytes if Necessary
func (e *encoder) breakUpInto8bit() {
	// fill ending code (max 4bit)
	// depends on max capacity of current version and EC level
	maxCap := e.version.NumTotalCodewords() * 8
	if less := maxCap - e.dst.Len(); less < 0 {
		err := fmt.Errorf(
			"wrong version(%d) cap(%d bits) and could not contain all bits: %d bits",
			e.version.Ver, maxCap, e.dst.Len(),
		)
		panic(err)
	} else if less < 4 {
		e.dst.AppendNumBools(less, false)
	} else {
		e.dst.AppendNumBools(4, false)
	}

	// append `0` to be 8 times bits length
	if mod := e.dst.Len() % 8; mod != 0 {
		e.dst.AppendNumBools(8-mod, false)
	}

	// _defaultPadding bytes
	// _defaultPadding byte 11101100 00010001
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
func (e *encoder) charCountBits() int {
	var lv int
	if v := e.version.Ver; v <= 9 {
		lv = 9
	} else if v <= 26 {
		lv = 26
	} else {
		lv = 40
	}
	pos := fmt.Sprintf("%d_%s", lv, getEncModeName(e.mode))
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
