package beautifier

import "github.com/yeqown/go-qrcode/v2"

type patternDesc interface {
	match(mat qrcode.Matrix) bool
}

type patternLoc struct {
}

type squarePattern struct {
	x, y int
}
