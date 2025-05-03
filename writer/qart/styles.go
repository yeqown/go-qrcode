package qart

import "github.com/yeqown/go-qrcode/v2"

type SVGPointExpr string

// Styler defines the style of QArt, how to response each block
// of QRCode.
type styler interface {
	Point(x, y int, v qrcode.QRValue) SVGPointExpr
}
