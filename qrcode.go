package qrcode

import (
	"github.com/yeqown/go-qrcode/version"
)

// New generate a QRCode struct to create or
func New(text string) (*QRCode, error) {
	qrc := &QRCode{
		content: text,
	}

	// initialize
	qrc.analyze()

	return qrc, nil
}

// QRCode contains: infos
// 1. data info
// 2. mask info
// etc.
type QRCode struct {
	content string            // input text content
	v       version.QRVersion // version means the size
}

func (q *QRCode) init() {
	q.analyze()
}

// Save QRCode image into saveToPath
func (q *QRCode) Save(saveToPath string) error {
	return nil
}

func (q *QRCode) analyze() {
	q.v = version.Analyze(q.content)
	q.v.Init()
}

func (q *QRCode) encode() {

}
