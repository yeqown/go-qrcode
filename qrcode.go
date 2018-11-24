package qrcode

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/skip2/go-qrcode/bitset"
	"github.com/skip2/go-qrcode/reedsolomon"
	"github.com/yeqown/go-qrcode/matrix"
)

var (
	// DEBUG mode flag
	DEBUG = false

	once sync.Once
)

// New generate a QRCode struct to create or
func New(text string) (*QRCode, error) {
	qrc := &QRCode{
		content: text,
	}

	// initialize
	if err := qrc.init(); err != nil {
		return nil, err
	}

	return qrc, nil
}

// QRCode contains: infos
type QRCode struct {
	content string // input text content
	rawData []byte // raw Data to transfer

	dataBSet *bitset.Bitset // final data bit stream of encode data
	mat      *matrix.Matrix // matrix grid to store final bitmap
	ecBSet   *bitset.Bitset // final error correction bitset

	v         Version  // version means the size
	ver       int      // version num
	recoverLv ECLevel  // recoveryLevel
	mode      EncMode  // EncMode
	encoder   *Encoder // encoder ptr to call it's methods ~
}

func (q *QRCode) init() error {
	once.Do(func() {
		if err := load(defaultVersionCfg); err != nil {
			panic(err)
		}
	})

	if err := q.analyze(); err != nil {
		return fmt.Errorf("could not analyze the data: %v", err)
	}

	q.rawData = []byte(q.content)
	q.mat = matrix.New(q.v.Dimension(), q.v.Dimension())
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

	// 交替排列
	q.arrarngeBits(dataBlocks, ecBlocks)

	// append ec after data
	q.dataBSet.Append(q.ecBSet)

	// append remainder bits
	q.dataBSet.AppendNumBools(q.v.RemainderBits, false)

	return nil
}

// analyze choose version and encoder
func (q *QRCode) analyze() error {
	// TODO: 选择版本
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
	var (
		bset *bitset.Bitset
	)
	bset, err = q.encoder.Encode(q.rawData)
	if err != nil {
		err = fmt.Errorf("could not encode data: %v", err)
		return
	}

	blocks = make([]dataBlock, q.v.TotalNumBlocks())

	// split bset into data Block
	start, end, blockID := 0, 0, 0
	for _, g := range q.v.Groups {
		for j := 0; j < g.NumBlocks; j++ {
			start = end
			end = start + g.NumDataCodewords*8

			blocks[blockID].Data = bset.Substr(start, end)
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
		debugLogf("numOfECBlock: %d", b.NumECBlock)
		bset := reedsolomon.Encode(b.Data, b.NumECBlock)
		blocks[idx].Data = bset.Substr(b.StartOffset, bset.Len())
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

// arrarngeBits ... and save into dataBSet
func (q *QRCode) arrarngeBits(dataBlocks []dataBlock, ecBlocks []ecBlock) {
	if DEBUG {
		log.Println("before arrange")
		for i := 0; i < len(ecBlocks); i++ {
			debugLogf("ec block_%d: %v", i, ecBlocks[i])
		}

		for i := 0; i < len(dataBlocks); i++ {
			debugLogf("data block_%d: %v", i, dataBlocks[i])
		}
	}
	// arrange data blocks
	var (
		overflowCnt = 0
		endFlag     = false
		curIdx      = 0
	)

	// check if bitsets initialized, or initial them
	if q.dataBSet == nil {
		q.dataBSet = bitset.New()
	}
	if q.ecBSet == nil {
		q.ecBSet = bitset.New()
	}

	for !endFlag {
		for _, block := range dataBlocks {
			start := curIdx * 8
			end := start + 8

			// debugLogf("arrange data blocks info: start: %d, end: %d, len: %d, overflowCnt: %d, curIdx: %d",
			// 	start, end, block.Data.Len(), overflowCnt, curIdx,
			// )

			if start >= block.Data.Len() {
				overflowCnt++
				continue
			}
			q.dataBSet.Append(block.Data.Substr(start, end))
		}
		curIdx = curIdx + 1

		if overflowCnt >= len(dataBlocks) {
			endFlag = true
		}
	}

	// arrange ec blocks, and reinitial
	endFlag = false
	overflowCnt = 0
	curIdx = 0

	for !endFlag {
		for _, block := range ecBlocks {
			start := curIdx * 8
			end := start + 8

			if start >= block.Data.Len() {
				overflowCnt++
				continue
			}
			q.ecBSet.Append(block.Data.Substr(start, end))
		}
		curIdx++

		if overflowCnt >= len(ecBlocks) {
			endFlag = true
		}
	}

	if DEBUG {
		if DEBUG {
			log.Println("after arrange")
			debugLogf("data bitsets: %s", q.dataBSet.String())
			debugLogf("ec bitsets: %s", q.ecBSet.String())
		}
	}
}

// InitMatrix with version info: ref to:
// http://www.thonky.com/qr-code-tutorial/module-placement-matrix
func (q *QRCode) initMatrix() {
	dimension := q.v.Dimension()
	if q.mat == nil {
		q.mat = matrix.New(dimension, dimension)
	}

	// add finder left-top
	addFinder(q.mat, 0, 0)
	addSpliter(q.mat, 7, 7, dimension)
	debugLogf("finish left-top finder")
	// add finder right-top
	addFinder(q.mat, dimension-7, 0)
	addSpliter(q.mat, dimension-8, 7, dimension)
	debugLogf("finish right-top finder")
	// add finder left-bottom
	addFinder(q.mat, 0, dimension-7)
	addSpliter(q.mat, 7, dimension-8, dimension)
	debugLogf("finish left-bottom finder")

	// 版本大于1
	if q.v.Ver > 1 {
		// add align-mode related to version cfg
		for _, loc := range loadAlignmentPatternLoc(q.v.Ver) {
			addAlign(q.mat, loc.X, loc.Y)
		}
		debugLogf("finish align")
	}
	// add timing line
	addTimingLine(q.mat, dimension)
	// add darkBlock always be position (4*ver+9, 8)
	addDarkBlock(q.mat, 8, 4*q.v.Ver+9)
	// presistFormatBlock for version and format info
	presistFormatBlock(q.mat, dimension)

	// presistVersionBlock for version over 7
	if q.v.Ver >= 7 {
		presistVersionBlock(q.mat, dimension)
	}
}

func addFinder(m *matrix.Matrix, top, left int) {
	// black outer
	x, y := top, left
	for i := 0; i < 24; i++ {
		m.Set(x, y, matrix.StateTrue)
		if i < 6 {
			x = x + 1
		} else if i < 12 {
			y = y + 1
		} else if i < 18 {
			x = x - 1
		} else {
			y = y - 1
		}
	}

	// white inner
	x, y = top+1, left+1
	for i := 0; i < 16; i++ {
		m.Set(x, y, matrix.StateFalse)
		if i < 4 {
			x = x + 1
		} else if i < 8 {
			y = y + 1
		} else if i < 12 {
			x = x - 1
		} else {
			y = y - 1
		}
	}

	// black inner
	for x := left + 2; x < left+5; x++ {
		for y := top + 2; y < top+5; y++ {
			m.Set(x, y, matrix.StateTrue)
		}
	}
}

func addSpliter(m *matrix.Matrix, x, y, dimension int) {
	// top-left
	if x == 7 && y == 7 {
		for pos := 0; pos < 8; pos++ {
			m.Set(x, pos, matrix.StateFalse)
			m.Set(pos, y, matrix.StateFalse)
		}
		return
	}

	// top-ritgh
	if x == dimension-8 && y == 7 {
		for pos := 0; pos < 8; pos++ {
			m.Set(x, y-pos, matrix.StateFalse)
			m.Set(x+pos, y, matrix.StateFalse)
		}
		return
	}

	// bottom-left
	if x == 7 && y == dimension-8 {
		for pos := 0; pos < 8; pos++ {
			m.Set(x, y+pos, matrix.StateFalse)
			m.Set(x-pos, y, matrix.StateFalse)
		}
		return
	}

}

// add matrix align
func addAlign(m *matrix.Matrix, centerX, centerY int) {
	m.Set(centerX, centerY, matrix.StateTrue)
	// black
	x, y := centerX-2, centerY-2
	for i := 0; i < 16; i++ {
		m.Set(x, y, matrix.StateTrue)
		if i < 4 {
			x = x + 1
		} else if i < 8 {
			y = y + 1
		} else if i < 12 {
			x = x - 1
		} else {
			y = y - 1
		}
	}
	// white
	x, y = centerX-1, centerY-1
	for i := 0; i < 8; i++ {
		m.Set(x, y, matrix.StateFalse)
		if i < 2 {
			x = x + 1
		} else if i < 4 {
			y = y + 1
		} else if i < 6 {
			x = x - 1
		} else {
			y = y - 1
		}
	}
}

// addTimingLine ...
func addTimingLine(m *matrix.Matrix, dimension int) {
	for pos := 8; pos < dimension-8; pos++ {
		if pos%2 == 0 {
			m.Set(6, pos, matrix.StateTrue)
			m.Set(pos, 6, matrix.StateTrue)
		} else {
			m.Set(6, pos, matrix.StateFalse)
			m.Set(pos, 6, matrix.StateFalse)
		}
	}
}

// addDarkBlock ...
func addDarkBlock(m *matrix.Matrix, x, y int) {
	m.Set(x, y, matrix.StateTrue)
}

// 为格式化信息保留模块
func presistFormatBlock(m *matrix.Matrix, dimension int) {

	// format-info
	for pos := 0; pos < 9; pos++ {
		// skip timing line
		if pos == 6 {
			m.Set(8, dimension-pos, matrix.StateFormat)
			m.Set(dimension-pos, 8, matrix.StateFormat)
			continue
		}
		// skip dark block
		if pos == 8 {
			m.Set(8, pos, matrix.StateFormat)           // top-left-column
			m.Set(pos, 8, matrix.StateFormat)           // top-left-row
			m.Set(dimension-pos, 8, matrix.StateFormat) // top-right-row
			continue
		}
		m.Set(8, pos, matrix.StateFormat)           // top-left-column
		m.Set(pos, 8, matrix.StateFormat)           // top-left-row
		m.Set(dimension-pos, 8, matrix.StateFormat) // top-right-row
		m.Set(8, dimension-pos, matrix.StateFormat) // bottom-left-column
	}

}

// 为版本信息保留模块
func presistVersionBlock(m *matrix.Matrix, dimension int) {
	// version info
	for i := 1; i <= 3; i++ {
		for pos := 0; pos < 6; pos++ {
			m.Set(dimension-8-i, pos, matrix.StateVersion)
			m.Set(pos, dimension-8-i, matrix.StateVersion)
		}
	}
}

// fillIntoMatrix fill q.dataBSet bitset stream into q.mat, ref to:
// http://www.thonky.com/qr-code-tutorial/module-placement-matrix
func (q *QRCode) fillIntoMatrix(dimension int) {
	var (
		x, y      = dimension - 1, dimension - 1
		l         = q.dataBSet.Len()
		upForward = true
		mod2, pos int

		setState, state matrix.State
		// turn      = false // if last loop, changed forward, this is true
		// downFoward = false
		// once sync.Once
		err error
	)

	for i := 0; pos < l; i++ {
		// debugLogf("fillIntoMatrix: dimension: %d, len: %d: pos: %d", dimension, l, pos)

		state, err = q.mat.Get(x, y)
		if err == matrix.ErrorOutRangeOfW {
			break
		}

		if q.dataBSet.At(pos) {
			setState = matrix.StateTrue
		} else {
			setState = matrix.StateFalse
		}

		if state == matrix.StateInit {
			q.mat.Set(x, y, setState)
			pos++
			debugLogf("normal set turn forward: upForward: %v, x: %d, y: %d", upForward, x, y)
		} else if state == matrix.ZERO {
			// turn forward and the new forward's block fisrt pos as value
			if upForward {
				x = x - 2
				y = y + 1
			} else {
				x = x - 2
				y = y - 1
			}

			if x == 7 || x == 6 {
				x = x - 1
			}

			upForward = !upForward
			debugLogf("unnormal state turn forward: upForward: %v, x: %d, y: %d", upForward, x, y)
			if s, _ := q.mat.Get(x, y); s == matrix.StateInit {
				q.mat.Set(x, y, setState)
				pos++
			}
		}

		// DO NOT CHANGE FOLLOWING CODE FOR NOW !!!
		// change x, y
		mod2 = i % 2

		// in one 8bit block
		if upForward {
			if mod2 == 0 {
				x = x - 1
			} else {
				y = y - 1
				x = x + 1
			}
		} else {
			if mod2 == 0 {
				x = x - 1
			} else {
				y = y + 1
				x = x + 1
			}
		}
	}

	debugLogf("fillDone and x: %d, y: %d, pos: %d, total: %d", x, y, pos, l)
}

// Save QRCode image into saveToPath
func (q *QRCode) Save(saveToPath string) error {
	// TODO: valid  saveToPath
	q.Draw()
	return drawAndSaveToFile(saveToPath, *q.mat)
}

// Draw ... Draw with bitset
func (q *QRCode) Draw() {
	dimension := q.v.Dimension()

	// 初始化二位矩阵
	q.initMatrix()

	// save current q.matrix copy
	mask0 := NewMask(q.mat, Modulo0)

	// fill bitset into matrix
	q.fillIntoMatrix(dimension)

	// draw("./testdata/bf_qrcode.jpeg", *q.mat)
	// XOR with mask and q.mat
	q.xorMask(mask0)
	// draw("./testdata/af_qrcode.jpeg", *q.mat)

	// fill format info
	q.fillFormatInfo(Modulo0, dimension)

	if q.v.Ver >= 7 {
		q.fillVersionInfo(dimension)
	}
}

// all mask patter and check the score choose the the lowest mask result
func (q *QRCode) xorMask(mask *Mask) {
	mask.mat.Iter(matrix.ROW, func(x, y int, s matrix.State) {
		// skip the empty palce
		if s == matrix.StateInit {
			return
		}
		s0, _ := q.mat.Get(x, y)
		q.mat.Set(x, y, matrix.XOR(s0, s))
	})
}

// fillVersionInfo ref to:
// https://www.thonky.com/qr-code-tutorial/format-version-tables
func (q *QRCode) fillVersionInfo(dimension int) {
	verBSet := q.v.verInfo()
	var mod3, mod6 int
	for pos := 0; pos < 18; pos++ {
		mod3 = pos % 3
		mod6 = pos % 6

		if verBSet.At(pos) {
			q.mat.Set(mod6, dimension-12+mod3, matrix.StateTrue)
			q.mat.Set(dimension-12+mod3, mod6, matrix.StateTrue)
		} else {
			q.mat.Set(mod6, dimension-12+mod3, matrix.StateFalse)
			q.mat.Set(dimension-12+mod3, mod6, matrix.StateTrue)
		}
	}
}

// fill format info ref to:
// https://www.thonky.com/qr-code-tutorial/format-version-tables
func (q *QRCode) fillFormatInfo(mode MaskPatternModulo, dimension int) {
	fmtBSet := q.v.formatInfo(int(mode))
	debugLogf("fmtBitSet: %s", fmtBSet.String())
	var (
		x, y = 0, dimension - 1
	)

	for pos := 0; pos < 15; pos++ {
		if fmtBSet.At(pos) {
			// row
			q.mat.Set(x, 8, matrix.StateTrue)
			// column
			q.mat.Set(8, y, matrix.StateTrue)
		} else {
			// row
			q.mat.Set(x, 8, matrix.StateFalse)
			// column
			q.mat.Set(8, y, matrix.StateFalse)
		}

		x = x + 1
		y = y - 1

		// row skip
		if x == 6 {
			x = 7
		} else if x == 8 {
			x = dimension - 8
		}

		// column skip
		if y == dimension-8 {
			y = 8
		} else if y == 6 {
			y = 5
		}
	}
}

func (q *QRCode) debugDraw() {
	if !DEBUG {
		return
	}
	drawAndSaveToFile("./testdata/qrtest_loop.jpeg", *q.mat)
	time.Sleep(time.Millisecond * 100)
}

func debugLogf(fmt string, v ...interface{}) {
	if !DEBUG {
		return
	}
	log.Printf(fmt, v...)
}
