package qrcode

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/yeqown/go-qrcode/matrix"
)

// Draw image with matrix

var (
	_defaultFilename = "default.jpeg"
	_defaultPadding  = 40
)

// drawAndSaveToFile image with matrix
func drawAndSaveToFile(name string, m matrix.Matrix, opt *outputImageOptions) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer f.Close()

	return drawAndSave(f, m, opt)
}

// drawAndSave save image into io.Writer
func drawAndSave(w io.Writer, m matrix.Matrix, opt *outputImageOptions) (err error) {
	if opt == nil {
		opt = _defaultOutputOption
	}

	img := draw(m, opt)

	if err = jpeg.Encode(w, img, nil); err != nil {
		err = fmt.Errorf("jpeg.Encode got err: %v", err)
	}

	return
}

// draw deal QRCode's matrix to be a image.Image
func draw(mat matrix.Matrix, opt *outputImageOptions) image.Image {
	// w as image width, h as image height
	w := mat.Width()*opt.qrBlockWidth() + 2*_defaultPadding
	h := w
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	// top-bottom _defaultPadding
	for posX := 0; posX < w; posX++ {
		for posY := 0; posY < _defaultPadding; posY++ {
			rgba.Set(posX, posY, opt.backgroundColor())
		}

		for posY := h - _defaultPadding; posY < h; posY++ {
			rgba.Set(posX, posY, opt.backgroundColor())
		}
	}

	// left-right _defaultPadding
	for posY := _defaultPadding; posY < h-_defaultPadding; posY++ {
		for posX := 0; posX < _defaultPadding; posX++ {
			rgba.Set(posX, posY, opt.backgroundColor())
		}

		for posX := w - _defaultPadding; posX < w; posX++ {
			rgba.Set(posX, posY, opt.backgroundColor())
		}
	}

	// iterate the matrix to draw each pixel
	mat.Iterate(matrix.ROW, func(x int, y int, v matrix.State) {
		xStart := x*opt.qrBlockWidth() + _defaultPadding
		yStart := y*opt.qrBlockWidth() + _defaultPadding
		xEnd := (x+1)*opt.qrBlockWidth() + _defaultPadding
		yEnd := (y+1)*opt.qrBlockWidth() + _defaultPadding

		// draw the block
		// TODO(@yeqown): make this abstract to Shape
		for posX := xStart; posX < xEnd; posX++ {
			for posY := yStart; posY < yEnd; posY++ {
				rgba.Set(posX, posY, opt.stateRGBA(v))
			}
		}
	})

	// DONE(@yeqown): add logo image
	if opt.logoImage() != nil {
		// draw logo image into rgba
		bound := opt.logo.Bounds()
		upperLeft, lowerRight := bound.Min, bound.Max
		logoWidth, logoHeight := lowerRight.X-upperLeft.X, lowerRight.Y-upperLeft.Y

		if !validLogoImage(w, h, logoWidth, logoHeight) {
			log.Printf("w=%d, h=%d, logoW=%d, logoH=%d, logo is over than 1/5 of QRCode \n",
				w, h, logoWidth, logoHeight)
			goto done
		}

		// DONE(@yeqown): calculate the xOffset and yOffset
		// which point(xOffset, yOffset) should icon upper-left to start
		xOffset, yOffset := (w-logoWidth)/2, (h-logoHeight)/2

		for posX := upperLeft.X; posX < lowerRight.X; posX++ {
			for posY := upperLeft.Y; posY < lowerRight.Y; posY++ {
				rgba.Set(posX+xOffset, posY+yOffset, opt.logo.At(posX, posY))
			}
		}
	}
done:
	return rgba
}

func validLogoImage(qrWidth, qrHeight, logoWidth, logoHeight int) bool {
	return qrWidth >= 5*logoWidth && qrHeight >= 5*logoHeight
}
