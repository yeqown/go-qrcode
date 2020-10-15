package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"

	"github.com/yeqown/go-qrcode/matrix"
)

// Draw image with matrix

var (
	defaultExpandPixel = 20
	defaultFilename    = "default.jpeg"
	padding            = 40
)

//// SetExpandPixel set defaultExpandPixel, default is 40
//func SetExpandPixel(n int) {
//	if n < 0 {
//		panic("could not set the negative integer")
//	}
//	defaultExpandPixel = n
//}

// drawAndSaveToFile image with matrix
func drawAndSaveToFile(name string, m matrix.Matrix) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer f.Close()

	return drawAndSave(f, m)
}

// drawAndSave save image into io.Writer
func drawAndSave(w io.Writer, m matrix.Matrix) error {
	img := draw(m, _defaultOutputOption)
	return save(w, img)
}

// draw deal qrcode matrix as a image.Image
func draw(mat matrix.Matrix, opt *outputImageOptions) image.Image {
	_stateToRGBA[matrix.StateFalse] = opt.backgroundColor
	_stateToRGBA[matrix.StateTrue] = opt.qrColor

	// w as image width, h as image height
	w := mat.Width()*defaultExpandPixel + 2*padding
	h := w

	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	// top-bottom padding
	for posX := 0; posX < w; posX++ {
		for posY := 0; posY < padding; posY++ {
			rgba.Set(posX, posY, opt.backgroundColor)
		}

		for posY := h - padding; posY < h; posY++ {
			rgba.Set(posX, posY, opt.backgroundColor)
		}
	}

	// left-right padding
	for posY := padding; posY < h-padding; posY++ {
		for posX := 0; posX < padding; posX++ {
			rgba.Set(posX, posY, opt.backgroundColor)
		}

		for posX := w - padding; posX < w; posX++ {
			rgba.Set(posX, posY, opt.backgroundColor)
		}
	}

	// iterate the matrix to draw each pixel
	mat.Iterate(matrix.ROW, func(x int, y int, v matrix.State) {
		xStart := x*defaultExpandPixel + padding
		yStart := y*defaultExpandPixel + padding
		xEnd := (x+1)*defaultExpandPixel + padding
		yEnd := (y+1)*defaultExpandPixel + padding

		// true for black, false for white
		for posX := xStart; posX < xEnd; posX++ {
			for posY := yStart; posY < yEnd; posY++ {
				rgba.Set(posX, posY, stateRGBA(v))
			}
		}
	})

	// TODO: add icon image
	return rgba
}

// save to file
func save(w io.Writer, img image.Image) error {
	if err := jpeg.Encode(w, img, nil); err != nil {
		return fmt.Errorf("could not save image into file with err: %v", err)
	}

	return nil
}

var (
	// _stateToRGBA state map tp color.Gray16
	_stateToRGBA = map[matrix.State]color.Color{
		matrix.StateFalse: hexToRGBA("#1aa6b7"),
		matrix.StateTrue:  hexToRGBA("#01c5c4"),
		//matrix.StateInit:  hexToRGBA("#1aa6b7"),
		//matrix.StateVersion: hexToRGBA("#444444"),
		// matrix.StateFormat:  hexToRGBA("#555555"),
	}

	// _defaultStateColor default color of undefined matrix.State
	// it shouldn't be used.
	_defaultStateColor = hexToRGBA("#ff414d")
)

func hexToRGBA(s string) color.RGBA {
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

	return c
}

// stateRGBA get color.Color by value State
func stateRGBA(v matrix.State) color.Color {
	if v, ok := _stateToRGBA[v]; ok {
		return v
	}

	return _defaultStateColor
}
