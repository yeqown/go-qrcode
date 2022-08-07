package standard

import (
	"fmt"
	"image/color"

	"github.com/yeqown/go-qrcode/v2"
)

type Color interface {
	color.Color

	// SetTransparent set color be transparent. if the color implementation doesn't support,
	// it should have no any side effect.
	SetTransparent()
}

// colorWrapper uniform colors usage in qrcode generation.
type colorWrapper struct {
	color.Color
}

func (w colorWrapper) RGBA() (r, g, b, a uint32) {
	return w.Color.RGBA()
}

func (w colorWrapper) SetTransparent() {
	switch w.Color.(type) {
	case color.RGBA:
		if rgba, ok := w.Color.(color.RGBA); ok {
			(&rgba).A = 0x00
		}
	case CompressedColor:
		if cc, ok := w.Color.(CompressedColor); ok {
			(&cc).A = 0x00
		}
	}
}

// CompressedColor use only uint16 bits to contains color information, since
// color.RGBA has 4 fields(uint8) to contains R\G\B\A, 32bits total.
type CompressedColor struct {
	Y uint8
	A uint8
}

func (c CompressedColor) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	a = uint32(c.A)
	a |= a << 8
	return y, y, y, a
}

var (
	// color_WHITE default acts by color.RGBA
	color_WHITE = color_WHITE_rgba
	// color_BLACK default acts by color.RGBA
	color_BLACK = color_BLACK_rgba

	color_WHITE_rgba Color = colorWrapper{Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}}
	color_BLACK_rgba Color = colorWrapper{Color: color.RGBA{R: 0, G: 0, B: 0, A: 255}}
	//color_WHITE_compressed Color = colorWrapper{Color: CompressedColor{Y: 255, A: 255}}
	//color_BLACK_compressed Color = colorWrapper{Color: CompressedColor{Y: 0, A: 255}}

	color_WHITE_compressed Color = colorWrapper{Color: color.Gray{Y: 255}}
	color_BLACK_compressed Color = colorWrapper{Color: color.Gray{Y: 0}}
)

var (
	// _STATE_MAPPING mapping matrix.State to color.RGBA in debug mode.
	_STATE_MAPPING = map[qrcode.QRType]Color{
		qrcode.QRType_INIT:     parseFromHex("#ffffff"), // [bg]
		qrcode.QRType_DATA:     parseFromHex("#cdc9c3"), // [bg]
		qrcode.QRType_VERSION:  parseFromHex("#000000"), // [fg]
		qrcode.QRType_FORMAT:   parseFromHex("#444444"), // [fg]
		qrcode.QRType_FINDER:   parseFromHex("#555555"), // [fg]
		qrcode.QRType_DARK:     parseFromHex("#2BA859"), // [fg]
		qrcode.QRType_SPLITTER: parseFromHex("#2BA859"), // [fg]
		qrcode.QRType_TIMING:   parseFromHex("#000000"), // [fg]
	}
)

// parseFromHex convert hex string into color.RGBA
func parseFromHex(s string) Color {
	c := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0xff,
	}

	var err error
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	if err != nil {
		panic(err)
	}

	return colorWrapper{Color: c}
}

func parseFromColor(c color.Color) Color {
	//rgba, ok := c.(color.RGBA)
	//if ok {
	//	return colorWrapper{Color: rgba}
	//}

	//_, _, _, a := c.RGBA()
	//return color.RGBA{
	//	R: uint8(r),
	//	G: uint8(g),
	//	B: uint8(b),
	//	A: uint8(a),
	//}

	return colorWrapper{Color: c}
}
