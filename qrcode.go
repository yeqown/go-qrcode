package qrcode

import (
	"fmt"

	"github.com/skip2/go-qrcode/bitset"
	"github.com/yeqown/go-qrcode/matrix"
	"github.com/yeqown/go-qrcode/version"

	"github.com/skip2/go-qrcode/reedsolomon"
)

// New generate a QRCode struct to create or
func New(text string) (*QRCode, error) {
	qrc := &QRCode{
		content: text,
	}

	// initialize
	qrc.init()

	return qrc, nil
}

// QRCode contains: infos
// 1. data info
// 2. mask info
// etc.
type QRCode struct {
	content   string                // input text content
	rawData   []byte                // raw Data to transfer
	stream    *bitset.Bitset        // bit stream of encode data
	mat       *matrix.Matrix        // matrix grid to store final bitmap
	v         version.QRVersion     // version means the size
	recoverLv version.RecoveryLevel // recoveryLevel
	encoder   Encoder
}

func (q *QRCode) init() error {
	if err := q.analyze(); err != nil {
		return fmt.Errorf("could not analyze the data: %v", err)
	}
	// version initial
	q.v.Init(q.recoverLv)

	q.rawData = []byte(q.content)
	q.mat = matrix.NewMatrix(q.v.Dimension(), q.v.Dimension())
	// TODO: choose encoder by what?
	q.encoder = Encoder{
		version: q.v,
		mode:    chooseMode(q.rawData), // default choose this mode
	}

	var err error
	q.stream, err = q.encoder.Encode(q.rawData)

	q.encode()
	if err != nil {
		return fmt.Errorf("could not encode the data: %v", err)
	}
	return nil
}

// analyze choose version and encoder
func (q *QRCode) analyze() error {
	q.v = version.Analyze(q.content)
	return nil
}

// Save QRCode image into saveToPath
func (q *QRCode) Save(saveToPath string) error {
	// TODO: valid  saveToPath
	return q.draw(saveToPath)
}

func (q *QRCode) draw(saveToPath string) error {
	return draw(saveToPath, *q.mat)
}

func (q *QRCode) encode() {
	reedsolomon.Encode()
}
