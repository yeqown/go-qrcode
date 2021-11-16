package qrcode

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/yeqown/reedsolomon"
	"github.com/yeqown/reedsolomon/binary"

	"github.com/yeqown/go-qrcode/v2/matrix"
)

var (
	// _debug mode flag
	_debug = false
)

// New generate a QRCode struct to create
func New(text string) (*QRCode, error) {
	dst := DefaultEncodingOption()
	return build(text, *dst)
}

// NewWith generate a QRCode struct with
// specified `ver`(QR version) and `ecLv`(Error Correction level)
func NewWith(text string, opts ...EncodeOption) (*QRCode, error) {
	dst := DefaultEncodingOption()
	for _, opt := range opts {
		opt.apply(dst)
	}

	return build(text, *dst)
}

func build(text string, option encodingOption) (*QRCode, error) {
	qrc := &QRCode{
		content:     text,
		mode:        option.EncMode,
		ecLv:        option.EcLevel,
		needAnalyze: true,
	}
	// initialize QRCode instance
	if err := qrc.init(); err != nil {
		return nil, err
	}

	qrc.masking()

	return qrc, nil
}

// QRCode contains fields to generate QRCode matrix, outputImageOptions to Draw image,
// etc.
type QRCode struct {
	content string // input text content
	rawData []byte // raw Data to transfer

	dataBSet *binary.Binary // final data bit stream of encode data
	mat      *matrix.Matrix // matrix grid to store final bitmap
	ecBSet   *binary.Binary // final error correction bitset

	v       version  // version means the size
	ver     int      // version num
	ecLv    ecLevel  // error correction level
	mode    encMode  // encMode
	encoder *encoder // encoder ptr to call its methods ~

	needAnalyze bool // auto analyze form content or specified `mode, recoverLv, ver`
}

func (q *QRCode) Save(w Writer) error {
	if w == nil {
		w = nonWriter{}
	}

	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("[WARNNING] [go-qrcode] close writer failed: %v\n", err)
		}
	}()

	return w.Write(*q.mat)
}

func (q *QRCode) Dimension() int {
	if q.mat == nil {
		return 0
	}

	return q.mat.Width()
}

func (q *QRCode) init() error {
	//once.Do(func() {
	// once load versions config file into memory
	// if err := load(defaultVersionCfg); err != nil {
	// 	panic(err)
	// }
	//})
	q.rawData = []byte(q.content)
	if q.needAnalyze {
		// analyze the input data to choose to adapt version
		if err := q.analyze(); err != nil {
			return fmt.Errorf("could not analyze the data: %v", err)
		}
	}
	// or check need params

	if !q.needAnalyze {
		// choose version without auto analyze
		q.v = loadVersion(q.ver, q.ecLv)
	}

	q.mat = matrix.New(q.v.Dimension(), q.v.Dimension())
	q.encoder = &encoder{
		mode:    q.mode,
		ecLv:    q.ecLv,
		version: q.v,
	}

	var (
		dataBlocks []dataBlock // data encoding blocks
		ecBlocks   []ecBlock   // error correction blocks
		err        error       // global error var
	)

	// data encoding, and be split into blocks
	if dataBlocks, err = q.dataEncoding(); err != nil {
		return err
	}

	// generate er bitsets, and also be split into blocks
	if ecBlocks, err = q.errorCorrectionEncoding(dataBlocks); err != nil {
		return err
	}

	// arrange data blocks and EC blocks
	q.arrangeBits(dataBlocks, ecBlocks)

	// append ec bits after data bits
	q.dataBSet.Append(q.ecBSet)

	// append remainder bits
	q.dataBSet.AppendNumBools(q.v.RemainderBits, false)

	// initial the 2d matrix
	q.prefillMatrix()

	return nil
}

// analyze rawData abd based on AUTO settings choose version and encoder
func (q *QRCode) analyze() error {
	if q.mode == EncModeAuto {
		// choose encode mode (num, alpha num, byte, Japanese)
		q.mode = analyzeEncodeModeFromRaw(q.rawData)
	}

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

	// split bitset into data Block
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
	NumECBlock  int // error correction codewords num per data block
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

// arrangeBits ... and save into dataBSet
func (q *QRCode) arrangeBits(dataBlocks []dataBlock, ecBlocks []ecBlock) {
	if _debug {
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

	// arrange ec blocks and reinitialize
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

// prefillMatrix with version info: ref to:
// http://www.thonky.com/qr-code-tutorial/module-placement-matrix
func (q *QRCode) prefillMatrix() {
	dimension := q.v.Dimension()
	if q.mat == nil {
		q.mat = matrix.New(dimension, dimension)
	}

	// add finder left-top
	addFinder(q.mat, 0, 0)
	addSplitter(q.mat, 7, 7, dimension)
	debugLogf("finish left-top finder")
	// add finder right-top
	addFinder(q.mat, dimension-7, 0)
	addSplitter(q.mat, dimension-8, 7, dimension)
	debugLogf("finish right-top finder")
	// add finder left-bottom
	addFinder(q.mat, 0, dimension-7)
	addSplitter(q.mat, 7, dimension-8, dimension)
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
	// only version 7 and larger version should add version info
	if q.v.Ver >= 7 {
		reserveVersionBlock(q.mat, dimension)
	}
}

// add finder module
func addFinder(m *matrix.Matrix, top, left int) {
	// black outer
	x, y := top, left
	for i := 0; i < 24; i++ {
		_ = m.Set(x, y, matrix.StateFinder)
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
		_ = m.Set(x, y, matrix.StateFalse)
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
			_ = m.Set(x, y, matrix.StateFinder)
		}
	}
}

// add splitter module
func addSplitter(m *matrix.Matrix, x, y, dimension int) {
	// top-left
	if x == 7 && y == 7 {
		for pos := 0; pos < 8; pos++ {
			_ = m.Set(x, pos, matrix.StateFalse)
			_ = m.Set(pos, y, matrix.StateFalse)
		}
		return
	}

	// top-right
	if x == dimension-8 && y == 7 {
		for pos := 0; pos < 8; pos++ {
			_ = m.Set(x, y-pos, matrix.StateFalse)
			_ = m.Set(x+pos, y, matrix.StateFalse)
		}
		return
	}

	// bottom-left
	if x == 7 && y == dimension-8 {
		for pos := 0; pos < 8; pos++ {
			_ = m.Set(x, y+pos, matrix.StateFalse)
			_ = m.Set(x-pos, y, matrix.StateFalse)
		}
		return
	}

}

// add matrix align module
func addAlignment(m *matrix.Matrix, centerX, centerY int) {
	_ = m.Set(centerX, centerY, matrix.StateTrue)
	// black
	x, y := centerX-2, centerY-2
	for i := 0; i < 16; i++ {
		_ = m.Set(x, y, matrix.StateTrue)
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
		_ = m.Set(x, y, matrix.StateFalse)
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
			_ = m.Set(6, pos, matrix.StateTrue)
			_ = m.Set(pos, 6, matrix.StateTrue)
		} else {
			_ = m.Set(6, pos, matrix.StateFalse)
			_ = m.Set(pos, 6, matrix.StateFalse)
		}
	}
}

// addDarkBlock ...
func addDarkBlock(m *matrix.Matrix, x, y int) {
	_ = m.Set(x, y, matrix.StateTrue)
}

// reserveFormatBlock maintain the position in matrix for format info
func reserveFormatBlock(m *matrix.Matrix, dimension int) {
	for pos := 1; pos < 9; pos++ {
		// skip timing line
		if pos == 6 {
			_ = m.Set(8, dimension-pos, matrix.StateFormat)
			_ = m.Set(dimension-pos, 8, matrix.StateFormat)
			continue
		}
		// skip dark module
		if pos == 8 {
			_ = m.Set(8, pos, matrix.StateFormat)           // top-left-column
			_ = m.Set(pos, 8, matrix.StateFormat)           // top-left-row
			_ = m.Set(dimension-pos, 8, matrix.StateFormat) // top-right-row
			continue
		}
		_ = m.Set(8, pos, matrix.StateFormat)           // top-left-column
		_ = m.Set(pos, 8, matrix.StateFormat)           // top-left-row
		_ = m.Set(dimension-pos, 8, matrix.StateFormat) // top-right-row
		_ = m.Set(8, dimension-pos, matrix.StateFormat) // bottom-left-column
	}

	// fix(@yeqown): b4b5ae3 reduced two format reversed blocks on top-left-column and top-left-row.
	_ = m.Set(0, 8, matrix.StateFormat)
	_ = m.Set(8, 0, matrix.StateFormat)
}

// reserveVersionBlock maintain the position in matrix for version info
func reserveVersionBlock(m *matrix.Matrix, dimension int) {
	// 3x6=18 cells
	for i := 1; i <= 3; i++ {
		for pos := 0; pos < 6; pos++ {
			_ = m.Set(dimension-8-i, pos, matrix.StateVersion)
			_ = m.Set(pos, dimension-8-i, matrix.StateVersion)
		}
	}
}

// fillIntoMatrix fill q.dataBSet bitset stream into q.mat, ref to:
// http://www.thonky.com/qr-code-tutorial/module-placement-matrix
func (q *QRCode) fillIntoMatrix(m *matrix.Matrix, dimension int) {
	var (
		x, y      = dimension - 1, dimension - 1
		l         = q.dataBSet.Len()
		upForward = true
		mod2, pos int

		setState, state matrix.State
		// turn      = false // if last loop, changed forward, this is true
		// downForward = false
		// once sync.Once
		err error
	)

	for i := 0; pos < l; i++ {
		// debugLogf("fillIntoMatrix: dimension: %d, len: %d: pos: %d", dimension, l, pos)

		state, err = m.Get(x, y)
		if err == matrix.ErrorOutRangeOfW {
			break
		}

		if q.dataBSet.At(pos) {
			setState = matrix.StateTrue
		} else {
			setState = matrix.StateFalse
		}

		if state == matrix.StateInit {
			_ = m.Set(x, y, setState)
			pos++
			// debugLogf("normal set turn forward: upForward: %v, x: %d, y: %d", upForward, x, y)
		} else if state == matrix.ZERO {
			// turn forward and the new forward's block first pos as value
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
			// debugLogf("unmoral state turn forward: upForward: %v, x: %d, y: %d", upForward, x, y)
			if s, _ := m.Get(x, y); s == matrix.StateInit {
				_ = m.Set(x, y, setState)
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

// draw from bitset to matrix.Matrix, calculate all mask modula score,
// then decide which mask to use according to the mask's score (the lowest one).
func (q *QRCode) masking() {
	type maskScore struct {
		Score int
		Idx   int
	}

	var (
		masks       = make([]*mask, 8)
		mats        = make([]*matrix.Matrix, 8)
		lowScore    = math.MaxInt32
		markMatsIdx int
		scoreChan   = make(chan maskScore, 8)
		wg          sync.WaitGroup
	)

	dimension := q.v.Dimension()

	// init mask and mats
	for i := 0; i < 8; i++ {
		masks[i] = newMask(q.mat, maskPatternModulo(i))
		mats[i] = q.mat.Copy()
	}

	// generate 8 matrix with mask
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			// fill bitset into matrix
			q.fillIntoMatrix(mats[i], dimension)

			// debug output
			if _debug {
				//_ = standard.DebugDraw(fmt.Sprintf("draft/mats_%d.jpeg", i), *mats[i], nil)
				//_ = standard.DebugDraw(fmt.Sprintf("draft/mask_%d.jpeg", i), *masks[i].mat, nil)
			}

			// xor with mask
			q.xorMask(mats[i], masks[i])
			if _debug {
				//_ = standard.DebugDraw(fmt.Sprintf("draft/mats_mask_%d.jpeg", i), *mats[i], nil)
			}

			// fill format info
			q.fillFormatInfo(mats[i], maskPatternModulo(i), dimension)
			// version7 and larger version has version info
			if q.v.Ver >= 7 {
				q.fillVersionInfo(mats[i], dimension)
			}

			// calculate score and decide the lowest score and Draw
			score := calculateScore(mats[i])
			debugLogf("cur idx: %d, score: %d, current lowest: mats[%d]:%d", i, score, markMatsIdx, lowScore)
			scoreChan <- maskScore{
				Score: score,
				Idx:   i,
			}

			if _debug {
				//_ = standard.DebugDraw(fmt.Sprintf("draft/qrcode_mask_%d.jpeg", i), *mats[i], nil)
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

// all mask patter and check the maskScore choose the the lowest mask result
func (q *QRCode) xorMask(m *matrix.Matrix, mask *mask) {
	mask.mat.Iterate(matrix.ROW, func(x, y int, s matrix.State) {
		// skip the empty place
		if s == matrix.StateInit {
			return
		}
		s0, _ := m.Get(x, y)
		_ = m.Set(x, y, matrix.XOR(s0, s))
	})
}

// fillVersionInfo ref to:
// https://www.thonky.com/qr-code-tutorial/format-version-tables
func (q *QRCode) fillVersionInfo(m *matrix.Matrix, dimension int) {
	bin := q.v.verInfo()

	// from high bit to lowest
	pos := 0
	for j := 5; j >= 0; j-- {
		for i := 1; i <= 3; i++ {
			if bin.At(pos) {
				_ = m.Set(dimension-8-i, j, matrix.StateTrue)
				_ = m.Set(j, dimension-8-i, matrix.StateTrue)
			} else {
				_ = m.Set(dimension-8-i, j, matrix.StateFalse)
				_ = m.Set(j, dimension-8-i, matrix.StateFalse)
			}

			pos++
		}
	}
}

// fill format info ref to:
// https://www.thonky.com/qr-code-tutorial/format-version-tables
func (q *QRCode) fillFormatInfo(m *matrix.Matrix, mode maskPatternModulo, dimension int) {
	fmtBSet := q.v.formatInfo(int(mode))
	debugLogf("fmtBitSet: %s", fmtBSet.String())
	var (
		x, y = 0, dimension - 1
	)

	for pos := 0; pos < 15; pos++ {
		if fmtBSet.At(pos) {
			// row
			_ = m.Set(x, 8, matrix.StateTrue)
			// column
			_ = m.Set(8, y, matrix.StateTrue)
		} else {
			// row
			_ = m.Set(x, 8, matrix.StateFalse)
			// column
			_ = m.Set(8, y, matrix.StateFalse)
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

// SetDebugMode open debug mode
func SetDebugMode() {
	_debug = true
}

func debugLogf(fmt string, v ...interface{}) {
	if !_debug {
		return
	}
	log.Printf(fmt, v...)
}
