package qrcode

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/yeqown/go-qrcode/matrix"
)

/* Draw image with matrix
 */

var (
	defaultExpandPixel = 20
	defaultFilename    = "default.jpeg"
	padding            = 20
)

// SetExpandPixel set defaultExpandPixel, default is 20
func SetExpandPixel(n int) {
	if n < 0 {
		panic("could not set the negative interger")
	}
	defaultExpandPixel = n
}

// draw image with matrix
func draw(name string, m matrix.Matrix) error {
	// w as image width, h as image height
	w, h := m.Width()*defaultExpandPixel, m.Height()*defaultExpandPixel
	// create file
	if len(name) == 0 {
		name = defaultFilename
	}

	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer f.Close()

	// draw into image
	var (
		gray16 = image.NewGray16(image.Rect(0, 0, w, h))
	)

	// TODO: add padding of the image
	//

	// iter the matrix to draw each pixel
	m.Iter(matrix.ROW, func(x int, y int, v matrix.State) {
		xStart := x * defaultExpandPixel
		yStart := y * defaultExpandPixel
		xEnd := (x + 1) * defaultExpandPixel
		yEnd := (y + 1) * defaultExpandPixel

		// true for black, false for white
		for posX := xStart; posX < xEnd; posX++ {
			for posY := yStart; posY < yEnd; posY++ {
				if posX == xStart || posY == yStart {
					gray16.SetGray16(posX, posY, matrix.LoadGray16(matrix.BORDER))
					continue
				}
				gray16.SetGray16(posX, posY, matrix.LoadGray16(v))
			}
		}
	})

	// save to file
	if err := jpeg.Encode(f, gray16, nil); err != nil {
		return fmt.Errorf("could not save image into file with err: %v", err)
	}
	return nil
}
