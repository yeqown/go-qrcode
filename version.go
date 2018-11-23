package qrcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/skip2/go-qrcode/bitset"
)

// ECLevel 错误回复级别定义
type ECLevel int

const (
	// Low :Level L: 7% error recovery.
	Low ECLevel = iota
	// Meddium :Level M: 15% error recovery. Good default choice.
	Meddium
	// Quart :Level Q: 25% error recovery.
	Quart
	// Highest :Level H: 30% error recovery.
	Highest

	// 格式化信息Bit位数
	formatInfoLengthBits = 15

	// 版本信息bit位数
	versionInfoLengthBits = 18
)

var (
	// 默认版本信息定义路径
	defaultVersionCfg     = "./versionCfg.json"
	errMissMatchedVersion = errors.New("could not match version! check the versionCfg.json file")
	versions              []Version
	// Each QR Code contains a 15-bit Format Information value.  The 15 bits
	// consist of 5 data bits concatenated with 10 error correction bits.
	//
	// The 5 data bits consist of:
	// - 2 bits for the error correction level (L=01, M=00, G=11, H=10).
	// - 3 bits for the data mask pattern identifier.
	//
	// formatBitSequence is a mapping from the 5 data bits to the completed 15-bit
	// Format Information value.
	//
	// For example, a QR Code using error correction level L, and data mask
	// pattern identifier 001:
	//
	// 01 | 001 = 01001 = 0x9
	// formatBitSequence[0x9].qrCode = 0x72f3 = 111001011110011
	formatBitSequence = []struct {
		regular uint32
		micro   uint32
	}{
		{0x5412, 0x4445},
		{0x5125, 0x4172},
		{0x5e7c, 0x4e2b},
		{0x5b4b, 0x4b1c},
		{0x45f9, 0x55ae},
		{0x40ce, 0x5099},
		{0x4f97, 0x5fc0},
		{0x4aa0, 0x5af7},
		{0x77c4, 0x6793},
		{0x72f3, 0x62a4},
		{0x7daa, 0x6dfd},
		{0x789d, 0x68ca},
		{0x662f, 0x7678},
		{0x6318, 0x734f},
		{0x6c41, 0x7c16},
		{0x6976, 0x7921},
		{0x1689, 0x06de},
		{0x13be, 0x03e9},
		{0x1ce7, 0x0cb0},
		{0x19d0, 0x0987},
		{0x0762, 0x1735},
		{0x0255, 0x1202},
		{0x0d0c, 0x1d5b},
		{0x083b, 0x186c},
		{0x355f, 0x2508},
		{0x3068, 0x203f},
		{0x3f31, 0x2f66},
		{0x3a06, 0x2a51},
		{0x24b4, 0x34e3},
		{0x2183, 0x31d4},
		{0x2eda, 0x3e8d},
		{0x2bed, 0x3bba},
	}

	// QR Codes version 7 and higher contain an 18-bit Version Information value,
	// consisting of a 6 data bits and 12 error correction bits.
	//
	// versionBitSequence is a mapping from QR Code version to the completed
	// 18-bit Version Information value.
	//
	// For example, a QR code of version 7:
	// versionBitSequence[0x7] = 0x07c94 = 000111110010010100
	versionBitSequence = []uint32{
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x07c94,
		0x085bc,
		0x09a99,
		0x0a4d3,
		0x0bbf6,
		0x0c762,
		0x0d847,
		0x0e60d,
		0x0f928,
		0x10b78,
		0x1145d,
		0x12a17,
		0x13532,
		0x149a6,
		0x15683,
		0x168c9,
		0x177ec,
		0x18ec4,
		0x191e1,
		0x1afab,
		0x1b08e,
		0x1cc1a,
		0x1d33f,
		0x1ed75,
		0x1f250,
		0x209d5,
		0x216f0,
		0x228ba,
		0x2379f,
		0x24b0b,
		0x2542e,
		0x26a64,
		0x27541,
		0x28c69,
	}
)

// load versions config into `versions`
func load(pathToCfg string) error {
	versions = make([]Version, 4*40)

	fd, err := os.OpenFile(pathToCfg, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}

	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	return json.Unmarshal(b, &versions)
}

type capacity struct {
	Numeric      int `json:"n"` // 版本对应的数字容量
	AlphaNumeric int `json:"a"` // 字符
	Byte         int `json:"b"` // 字节
	JP           int `json:"j"` // 日文
}

// 这些值可用于确定给定Reed-Solomon块需要多少数据字节和纠错字节。
type group struct {
	// 分块数 num of blocks
	NumBlocks int `json:"nbs"`

	// 数据代码字数 Number of data codewords.
	NumDataCodewords int `json:"ndcs"`

	// 每一个块的EC代码字数
	ECBlockwordsPerBlock int `json:"ecbs_pb"`
}

// Version ...
type Version struct {
	// 版本号 int 1-40
	Ver int `json:"ver"`

	// 错误恢复级别 0, 1, 2, 3
	ECLevel ECLevel `json:"eclv"`

	// ref to: https://www.thonky.com/qr-code-tutorial/character-capacities
	// 每一个版本 & 回复级别 对应的最大bit容量
	Cap capacity `json:"cap"`

	// RemainderBits 对应版本需要补充的剩余位bit数
	RemainderBits int `json:"rembits"`

	// Blocks 生成纠错码需要的分组分块信息
	// ref to: https://www.thonky.com/qr-code-tutorial/error-correction-table
	// numGroup = len(Groups)
	Groups []group `json:"groups"`
}

// Dimension ...
func (v Version) Dimension() int {
	return v.Ver*4 + 17
}

// NumTotalCodewrods 总的数据代码字数
func (v Version) NumTotalCodewrods() int {
	var total int
	for _, g := range v.Groups {
		total = total + (g.NumBlocks * g.NumDataCodewords)
	}
	return total
}

// NumGroups ... 分组数
func (v Version) NumGroups() int {
	return len(v.Groups)
}

// TotalNumBlocks ...总的分块数
func (v Version) TotalNumBlocks() int {
	var total int
	for _, g := range v.Groups {
		total = total + g.NumBlocks
	}
	return total
}

// VerInfo Version info bitset
func (v Version) verInfo() *bitset.Bitset {
	if v.Ver < 7 {
		return nil
	}

	result := bitset.New()
	result.AppendUint32(versionBitSequence[v.Ver], 18)

	return result
}

// formatInfo returns the 15-bit Format Information value for a QR
// code.
func (v Version) formatInfo(maskPattern int) *bitset.Bitset {
	formatID := 0

	switch v.ECLevel {
	case Low:
		formatID = 0x08 // 0b01000
	case Meddium:
		formatID = 0x00 // 0b00000
	case Quart:
		formatID = 0x18 // 0b11000
	case Highest:
		formatID = 0x10 // 0b10000
	default:
		log.Panicf("Invalid level %d", v.ECLevel)
	}

	if maskPattern < 0 || maskPattern > 7 {
		log.Panicf("Invalid maskPattern %d", maskPattern)
	}

	formatID |= maskPattern & 0x7
	result := bitset.New()
	result.AppendUint32(formatBitSequence[formatID].regular, formatInfoLengthBits)
	return result
}

// loadVersion get version config from config
func loadVersion(lv int, ecLv ECLevel) Version {
	for _, v := range versions {
		if v.Ver == lv && v.ECLevel == ecLv {
			return v
		}
	}
	panic(errMissMatchedVersion)
}

// Analyze the text, and decide which version should be choose
// ref to: http://muyuchengfeng.xyz/%E4%BA%8C%E7%BB%B4%E7%A0%81-%E5%AD%97%E7%AC%A6%E5%AE%B9%E9%87%8F%E8%A1%A8/
func Analyze(text string) Version {
	return loadVersion(1, Low)
}

// SetVersionCfgFile set custom version config file
func SetVersionCfgFile(fp string) {
	defaultVersionCfg = fp
}
