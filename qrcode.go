package qrcode

import (
	"fmt"
	"log"

	"github.com/skip2/go-qrcode/bitset"
	"github.com/skip2/go-qrcode/reedsolomon"
	"github.com/yeqown/go-qrcode/matrix"
)

var (
	// DEBUG mode flag
	DEBUG = true
)

// NewQRCode generate a QRCode struct to create or
func NewQRCode(text string) (*QRCode, error) {
	qrc := &QRCode{
		content: text,
	}

	// initialize
	qrc.init()

	return qrc, nil
}

// QRCode contains: infos
type QRCode struct {
	content string // input text content
	rawData []byte // raw Data to transfer

	dataBSet *bitset.Bitset // final data bit stream of encode data
	ecBSet   *bitset.Bitset // final error correction bitset
	mat      *matrix.Matrix // matrix grid to store final bitmap

	v         Version  // version means the size
	ver       int      // version num
	recoverLv ECLevel  // recoveryLevel
	mode      EncMode  // EncMode
	encoder   *Encoder // encoder ptr to call it's methods ~
}

func (q *QRCode) init() error {
	if err := q.analyze(); err != nil {
		return fmt.Errorf("could not analyze the data: %v", err)
	}

	q.rawData = []byte(q.content)
	q.mat = matrix.NewMatrix(q.v.Dimension(), q.v.Dimension())
	q.encoder = &Encoder{
		mode:    q.mode,
		ecLv:    q.recoverLv,
		version: q.v,
	}

	var (
		dataBlocks []dataBlock
		ecBlocks   []ecBlock
		err        error
	)

	// 数据编码
	if dataBlocks, err = q.dataEncoding(); err != nil {
		return err
	}

	// 生成纠错码
	if ecBlocks, err = q.errorCorrectionEncoding(dataBlocks); err != nil {
		return err
	}

	q.arrarngeBits(dataBlocks, ecBlocks)

	return nil
}

// analyze choose version and encoder
func (q *QRCode) analyze() error {
	// 选择版本
	q.ver = 5
	// 选择错误矫正级别
	q.recoverLv = Quart
	// 确定版本
	q.v = loadVersion(q.ver, q.recoverLv)
	// 确定模式
	q.mode = EncModeByte

	// TODO: analyze content to decide version and mode. etc.
	// q.v = Analyze(q.content)
	return nil
}

// dataEncoding ref to:
// https://www.thonky.com/qr-code-tutorial/data-encoding
func (q *QRCode) dataEncoding() (blocks []dataBlock, err error) {
	// var bset *bitset.Bitset
	blocks = make([]dataBlock, q.v.TotalNumBlocks())

	_, err = q.encoder.Encode(q.rawData)
	if err != nil {
		err = fmt.Errorf("could not encode data: %v", err)
		return
	}

	// split bset into data Block
	start, end, blockID := 0, 0, 0
	for _, g := range q.v.Groups {
		for j := 0; j < g.NumBlocks; j++ {
			start = end
			end = start + g.NumDataCodewords*8

			blocks[blockID].Data = q.dataBSet.Substr(start, end)
			blocks[blockID].StartOffset = end - start
			blocks[blockID].NumECBlock = g.ECBlockwordsPerBlock

			blockID++
		}
	}

	return
}

// dataBlock ...
type dataBlock struct {
	Data        *bitset.Bitset
	StartOffset int
	NumECBlock  int
}

// ecBlock ...
type ecBlock struct {
	Data        *bitset.Bitset
	StartOffset int
}

// errorCorrectionEncoding ref to:
// https://www.thonky.com/qr-code-tutorial /error-correction-coding
func (q *QRCode) errorCorrectionEncoding(dataBlocks []dataBlock) (blocks []ecBlock, err error) {
	// start, end, blockID := 0, 0, 0
	blocks = make([]ecBlock, q.v.TotalNumBlocks())

	for idx, b := range dataBlocks {
		blocks[idx].Data = reedsolomon.Encode(b.Data, b.NumECBlock)
		blocks[idx].StartOffset = b.StartOffset
	}

	return

	// 分组，分块
	// loop group
	// for _, g := range q.v.Groups {
	// 	// loop block
	// 	for j := 0; j < g.NumBlocks; j++ {
	// 		start = end
	// 		end = start + g.NumDataCodewords*8

	// 		blocks[blockID].Data = reedsolomon.Encode(q.dataBSet.Substr(start, end), g.ECBlockwordsPerBlock)
	// 		blocks[blockID].StartOffset = end - start

	// 		blockID++
	// 	}
	// }
	// return
}

// arrarngeBits ...
func (q *QRCode) arrarngeBits(dataBlocks []dataBlock, ecBlocks []ecBlock) {
	if DEBUG {
		for i := 0; i < len(ecBlocks); i++ {
			debugLogf("ec block_%d: %v", i, ecBlocks[i])
		}

		for i := 0; i < len(dataBlocks); i++ {
			debugLogf("data block_%d: %v", i, dataBlocks[i])
		}
	}
	// TODO: arrange data & EC blocks into bitsets
}

// Append remainder bits.
func (q *QRCode) appendRemainderBits() {

}

// Save QRCode image into saveToPath
func (q *QRCode) Save(saveToPath string) error {
	// TODO: valid  saveToPath
	return q.draw(saveToPath)
}

func (q *QRCode) draw(saveToPath string) error {
	return draw(saveToPath, *q.mat)
}

func debugLogf(fmt string, v ...interface{}) {
	if !DEBUG {
		return
	}
	log.Printf(fmt, v...)
}
