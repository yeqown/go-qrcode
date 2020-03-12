package qrcode

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
	"time"

	// "github.com/skip2/go-qrcode/bitset"
	// "github.com/skip2/go-qrcode/reedsolomon"

	"github.com/yeqown/go-qrcode/matrix"
	"github.com/yeqown/reedsolomon"
	"github.com/yeqown/reedsolomon/binary"
)

var (
	// DEBUG mode flag
	DEBUG = false

	// once to load versions config file
	once sync.Once
)

// New generate a QRCode struct to create
func New(text string) (*QRCode, error) {
	qrc := &QRCode{
		content:     text,
		mode:        EncModeByte,
		needAnalyze: true,
	}

	// initialize QRCode instance
	if err := qrc.init(); err != nil {
		return nil, err
	}

	return qrc, nil
}

// NewWithSpecV generate a QRCode struct with
// specified `ver`(QR version) and `ecLv`(Error Correction level)
func NewWithSpecV(text string, ver int, ecLv ECLevel) (*QRCode, error) {
	qrc := &QRCode{
		content:     text,
		ver:         ver,
		mode:        EncModeByte,
		ecLv:        ecLv,
		needAnalyze: false,
	}
	// initialize QRCode instance
	if err := qrc.init(); err != nil {
		return nil, err
	}

	return qrc, nil
}

// QRCode contains: infos
type QRCode struct {
	content string // input text content
	rawData []byte // raw Data to transfer

	dataBSet *binary.Binary // final data bit stream of encode data
	mat      *matrix.Matrix // matrix grid to store final bitmap
	ecBSet   *binary.Binary // final error correction bitset

	v       Version  // version means the size
	ver     int      // version num
	ecLv    ECLevel  // recoveryLevel
	mode    EncMode  // EncMode
	encoder *Encoder // encoder ptr to call it's methods ~

	needAnalyze bool // auto analyze form content or specified `mode, recoverLv, ver`
}

func (q *QRCode) init() error {
	once.Do(func() {
		// once load versions config file into memory
		// if err := load(defaultVersionCfg); err != nil {
		// 	panic(err)
		// }
	})
	q.rawData = []byte(q.content)
	if q.needAnalyze {
		// analyze the input data to choose adapt version
		if err := q.analyze(); err != nil {
			return fmt.Errorf("could not analyze the data: %v", err)
		}
	}
	// or check need params

	// choose version without auto analyze
	if !q.needAnalyze {
		q.v = loadVersion(q.ver, q.ecLv)
	}

	q.mat = matrix.New(q.v.Dimension(), q.v.Dimension())
	q.encoder = &Encoder{
		mode:    q.mode,
		ecLv:    q.ecLv,
		version: q.v,
	}

	var (
		dataBlocks []dataBlock // data encoding blocks
		ecBlocks   []ecBlock   // error correction blocks
		err        error       // global error var
	)

	// data encoding, and be splited into blocks
	if dataBlocks, err = q.dataEncoding(); err != nil {
		return err
	}

	// generate er bitsets, and alse be spilited into blocks
	if ecBlocks, err = q.errorCorrectionEncoding(dataBlocks); err != nil {
		return err
	}

	// arrange datablocks and ecblocks
	q.arrarngeBits(dataBlocks, ecBlocks)

	// append ec bits after data bits
	q.dataBSet.Append(q.ecBSet)

	// append remainder bits
	q.dataBSet.AppendNumBools(q.v.RemainderBits, false)

	return nil
}

// analyze choose version and encoder
func (q *QRCode) analyze() error {
	// choose error correction level
	q.ecLv = Quart

	// choose encode mode (num, alphanum, byte, Japanese)
	q.mode = anlayzeMode(q.rawData)

	// analyze content to decide version etc.
	analyzedV, err := analyzeVersion(q.rawData, q.ecLv, q.mode)
	if err != nil {
		return fmt.Errorf("could not analyzeVersion: %v", err)
	}
	q.v = *analyzedV
	q.ver = (*analyzedV).Ver
	return nil
}

// dataEncoding ref to:
// https://www.thonky.com/qr-code-tutorial/data-encoding
func (q *QRCode) dataEncoding() (blocks []dataBlock, err error) {
	var (
		bset *binary.Binary
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

			blocks[blockID].Data, err = bset.Subset(start, end)
			if err != nil {
				panic(err)
			}
			blocks[blockID].StartOffset = end - start
			blocks[blockID].NumECBlock = g.ECBlockwordsPerBlock

			blockID++
		}
	}

	return
}

// dataBlock ...
type dataBlock struct {
	Data        *binary.Binary
	StartOffset int // length
	NumECBlock  int // error correction codewrods num per data block
}

// ecBlock ...
type ecBlock struct {
	Data *binary.Binary
	// StartOffset int // length
}

// errorCorrectionEncoding ref to:
// https://www.thonky.com/qr-code-tutorial /error-correction-coding
func (q *QRCode) errorCorrectionEncoding(dataBlocks []dataBlock) (blocks []ecBlock, err error) {
	// start, end, blockID := 0, 0, 0
	blocks = make([]ecBlock, q.v.TotalNumBlocks())
	for idx, b := range dataBlocks {
		debugLogf("numOfECBlock: %d", b.NumECBlock)
		bset := reedsolomon.Encode(b.Data, b.NumECBlock)
		blocks[idx].Data, err = bset.Subset(b.StartOffset, bset.Len())
		if err != nil {
			panic(err)
		}
		// blocks[idx].StartOffset = b.StartOffset
	}
	return
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
		start, end  int
	)

	// check if bitsets initialized, or initial them
	if q.dataBSet == nil {
		q.dataBSet = binary.New()
	}
	if q.ecBSet == nil {
		q.ecBSet = binary.New()
	}

	for !endFlag {
		for _, block := range dataBlocks {
			start = curIdx * 8
			end = start + 8
			if start >= block.Data.Len() {
				overflowCnt++
				continue
			}
			subBin, err := block.Data.Subset(start, end)
			if err != nil {
				panic(err)
			}
			q.dataBSet.Append(subBin)
			debugLogf("arrange data blocks info: start: %d, end: %d, len: %d, overflowCnt: %d, curIdx: %d",
				start, end, block.Data.Len(), overflowCnt, curIdx,
			)
		}
		curIdx++
		// loop finish check
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
			start = curIdx * 8
			end = start + 8

			if start >= block.Data.Len() {
				overflowCnt++
				continue
			}
			subBin, err := block.Data.Subset(start, end)
			if err != nil {
				panic(err)
			}
			q.ecBSet.Append(subBin)
		}
		curIdx++
		// loop finish check
		if overflowCnt >= len(ecBlocks) {
			endFlag = true
		}
	}

	debugLogf("after arrange")
	debugLogf("data bitsets: %s", q.dataBSet.String())
	debugLogf("ec bitsets: %s", q.ecBSet.String())
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

	// only version-1 QR code has no alignment module
	if q.v.Ver > 1 {
		// add align-mode related to version cfg
		for _, loc := range loadAlignmentPatternLoc(q.v.Ver) {
			addAlignment(q.mat, loc.X, loc.Y)
		}
		debugLogf("finish align")
	}
	// add timing line
	addTimingLine(q.mat, dimension)
	// add darkBlock always be position (4*ver+9, 8)
	addDarkBlock(q.mat, 8, 4*q.v.Ver+9)
	// reserveFormatBlock for version and format info
	reserveFormatBlock(q.mat, dimension)

	// reserveVersionBlock for version over 7
	// only version 7 and larger version shoud add Version info
	if q.v.Ver >= 7 {
		reserveVersionBlock(q.mat, dimension)
	}
}

// add finder module
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

// add spliter module
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

// add matrix align module
func addAlignment(m *matrix.Matrix, centerX, centerY int) {
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

// reserveFormatBlock maitain the postion in matrix for format info
func reserveFormatBlock(m *matrix.Matrix, dimension int) {
	for pos := 0; pos < 9; pos++ {
		// skip timing line
		if pos == 6 {
			m.Set(8, dimension-pos, matrix.StateFormat)
			m.Set(dimension-pos, 8, matrix.StateFormat)
			continue
		}
		// skip dark module
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

// reserveVersionBlock maitain the postion in matrix for version info
func reserveVersionBlock(m *matrix.Matrix, dimension int) {
	// 3x6=18 cells
	for i := 1; i <= 3; i++ {
		for pos := 0; pos < 6; pos++ {
			m.Set(dimension-8-i, pos, matrix.StateVersion)
			m.Set(pos, dimension-8-i, matrix.StateVersion)
		}
	}
}

// fillIntoMatrix fill q.dataBSet bitset stream into q.mat, ref to:
// http://www.thonky.com/qr-code-tutorial/module-placement-matrix
func (q *QRCode) fillIntoMatrix(mat *matrix.Matrix, dimension int) {
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

		state, err = mat.Get(x, y)
		if err == matrix.ErrorOutRangeOfW {
			break
		}

		if q.dataBSet.At(pos) {
			setState = matrix.StateTrue
		} else {
			setState = matrix.StateFalse
		}

		if state == matrix.StateInit {
			mat.Set(x, y, setState)
			pos++
			// debugLogf("normal set turn forward: upForward: %v, x: %d, y: %d", upForward, x, y)
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
			// debugLogf("unnormal state turn forward: upForward: %v, x: %d, y: %d", upForward, x, y)
			if s, _ := mat.Get(x, y); s == matrix.StateInit {
				mat.Set(x, y, setState)
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
	if _, err := os.Open(saveToPath); err != nil && os.IsExist(err) {
		log.Printf("could not find path: %s, then save to %s",
			saveToPath, defaultFilename)
		saveToPath = defaultFilename
	}
	q.draw()
	return drawAndSaveToFile(saveToPath, *q.mat)
}

// SaveTo QRCode image into `w`(io.Writer)
func (q *QRCode) SaveTo(w io.Writer) error {
	q.draw()
	return drawAndSave(w, *q.mat)
}

// draw ... draw with bitset
func (q *QRCode) draw() {

	type sc struct {
		Score int
		Idx   int
	}
	var (
		masks       = make([]*Mask, 8)
		mats        = make([]*matrix.Matrix, 8)
		lowScore    = math.MaxInt32
		markMatsIdx int
		scoreChan   = make(chan sc, 8)
		wg          sync.WaitGroup
	)

	dimension := q.v.Dimension()

	// initial the 2d matrix
	q.initMatrix()

	// init mask and mats
	for i := 0; i < 8; i++ {
		masks[i] = NewMask(q.mat, MaskPatternModulo(i))
		mats[i] = q.mat.Copy()
	}

	// generate 8 matrix with mask
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			// fill bitset into matrix
			q.fillIntoMatrix(mats[i], dimension)

			// debug output
			if DEBUG {
				drawAndSaveToFile(fmt.Sprintf("draft/mats_%d.jpeg", i), *mats[i])
				drawAndSaveToFile(fmt.Sprintf("draft/mask_%d.jpeg", i), *masks[i].mat)
			}

			// xor with mask
			q.xorMask(mats[i], masks[i])
			if DEBUG {
				drawAndSaveToFile(fmt.Sprintf("draft/mats_mask_%d.jpeg", i), *mats[i])
			}

			// fill format info
			q.fillFormatInfo(mats[i], MaskPatternModulo(i), dimension)
			// verion7 and larger version has version info
			if q.v.Ver >= 7 {
				q.fillVersionInfo(mats[i], dimension)
			}

			// calculate score and decide the low score and draw
			score := CalculateScore(mats[i])
			debugLogf("cur idx: %d, score: %d, current lowest: mats[%d]:%d", i, score, markMatsIdx, lowScore)
			scoreChan <- sc{
				Score: score,
				Idx:   i,
			}
			// if score < lowScore {
			// 	lowScore = score
			// 	markMatsIdx = i
			// }
			if DEBUG {
				drawAndSaveToFile(fmt.Sprintf("draft/qrcode_mask_%d.jpeg", i), *mats[i])
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	close(scoreChan)

	for c := range scoreChan {
		if c.Score < lowScore {
			lowScore = c.Score
			markMatsIdx = c.Idx
		}
	}

	q.mat = mats[markMatsIdx]
}

// all mask patter and check the score choose the the lowest mask result
func (q *QRCode) xorMask(mat *matrix.Matrix, mask *Mask) {
	mask.mat.Iter(matrix.ROW, func(x, y int, s matrix.State) {
		// skip the empty palce
		if s == matrix.StateInit {
			return
		}
		s0, _ := mat.Get(x, y)
		mat.Set(x, y, matrix.XOR(s0, s))
	})
}

// fillVersionInfo ref to:
// https://www.thonky.com/qr-code-tutorial/format-version-tables
func (q *QRCode) fillVersionInfo(mat *matrix.Matrix, dimension int) {
	verBSet := q.v.verInfo()
	var mod3, mod6 int
	for pos := 0; pos < 18; pos++ {
		mod3 = pos % 3
		mod6 = pos % 6

		if verBSet.At(pos) {
			mat.Set(mod6, dimension-12+mod3, matrix.StateTrue)
			mat.Set(dimension-12+mod3, mod6, matrix.StateTrue)
		} else {
			mat.Set(mod6, dimension-12+mod3, matrix.StateFalse)
			mat.Set(dimension-12+mod3, mod6, matrix.StateTrue)
		}
	}
}

// fill format info ref to:
// https://www.thonky.com/qr-code-tutorial/format-version-tables
func (q *QRCode) fillFormatInfo(mat *matrix.Matrix, mode MaskPatternModulo, dimension int) {
	fmtBSet := q.v.formatInfo(int(mode))
	debugLogf("fmtBitSet: %s", fmtBSet.String())
	var (
		x, y = 0, dimension - 1
	)

	for pos := 0; pos < 15; pos++ {
		if fmtBSet.At(pos) {
			// row
			mat.Set(x, 8, matrix.StateTrue)
			// column
			mat.Set(8, y, matrix.StateTrue)
		} else {
			// row
			mat.Set(x, 8, matrix.StateFalse)
			// column
			mat.Set(8, y, matrix.StateFalse)
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
