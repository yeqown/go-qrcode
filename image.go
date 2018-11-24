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

/* Draw image with matrix
 */

var (
	defaultExpandPixel = 20
	defaultFilename    = "default.jpeg"
	padding            = 40
)

// SetExpandPixel set defaultExpandPixel, default is 40
func SetExpandPixel(n int) {
	if n < 0 {
		panic("could not set the negative interger")
	}
	defaultExpandPixel = n
}

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
	img := draw(m)
	return save(w, img)
}

func draw(mat matrix.Matrix) image.Image {
	// w as image width, h as image height
	w := mat.Width()*defaultExpandPixel + 2*padding
	h := w

	// draw into image
	var (
		gray16 = image.NewGray16(image.Rect(0, 0, w, h))
	)

	// top-bottom padding
	for posX := 0; posX < w; posX++ {
		for posY := 0; posY < padding; posY++ {
			gray16.SetGray16(posX, posY, color.White)
		}

		for posY := h - padding; posY < h; posY++ {
			gray16.SetGray16(posX, posY, color.White)
		}
	}

	// left-right padding
	for posY := padding; posY < h-padding; posY++ {
		for posX := 0; posX < padding; posX++ {
			gray16.SetGray16(posX, posY, color.White)
		}

		for posX := w - padding; posX < w; posX++ {
			gray16.SetGray16(posX, posY, color.White)
		}
	}

	// iter the matrix to draw each pixel
	mat.Iter(matrix.ROW, func(x int, y int, v matrix.State) {
		xStart := x*defaultExpandPixel + padding
		yStart := y*defaultExpandPixel + padding
		xEnd := (x+1)*defaultExpandPixel + padding
		yEnd := (y+1)*defaultExpandPixel + padding

		// true for black, false for white
		for posX := xStart; posX < xEnd; posX++ {
			for posY := yStart; posY < yEnd; posY++ {
				// block border
				// if posX == xStart || posY == yStart {
				// 	gray16.SetGray16(posX, posY, matrix.LoadGray16(matrix.BORDER))
				// 	continue
				// }
				gray16.SetGray16(posX, posY, matrix.LoadGray16(v))
			}
		}
	})

	return gray16
}

func save(w io.Writer, img image.Image) error {

	// save to file
	if err := jpeg.Encode(w, img, nil); err != nil {
		return fmt.Errorf("could not save image into file with err: %v", err)
	}

	return nil
}
