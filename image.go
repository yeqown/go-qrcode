package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"

	"github.com/fogleman/gg"

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

	var encoder ImageEncoder
	switch opt.formatTyp() {
	case JPEG_FORMAT:
		encoder = jpegEncoder{}
	case PNG_FORMAT:
		encoder = pngEncoder{}
	case HEIF_FORMAT:
		encoder = heifEncoder{}
	default:
		panic("Not supported file format")
	}

	// DONE(@yeqown): support file format specified config option
	if err = encoder.Encode(w, img); err != nil {
		err = fmt.Errorf("jpeg.Encode got err: %v", err)
	}

	return
}

//
//// draw deal QRCode's matrix to be a image.Image
//func draw(mat matrix.Matrix, opt *outputImageOptions) image.Image {
//	// w as image width, h as image height
//	w := mat.Width()*opt.qrBlockWidth() + 2*_defaultPadding
//	h := w
//	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
//
//	// top-bottom _defaultPadding
//	for posX := 0; posX < w; posX++ {
//		for posY := 0; posY < _defaultPadding; posY++ {
//			rgba.Set(posX, posY, opt.backgroundColor())
//		}
//
//		for posY := h - _defaultPadding; posY < h; posY++ {
//			rgba.Set(posX, posY, opt.backgroundColor())
//		}
//	}
//
//	// left-right _defaultPadding
//	for posY := _defaultPadding; posY < h-_defaultPadding; posY++ {
//		for posX := 0; posX < _defaultPadding; posX++ {
//			rgba.Set(posX, posY, opt.backgroundColor())
//		}
//
//		for posX := w - _defaultPadding; posX < w; posX++ {
//			rgba.Set(posX, posY, opt.backgroundColor())
//		}
//	}
//
//	ctx := &DrawContext{
//		upperLeft:  image.Point{}, // useless
//		lowerRight: image.Point{}, // useless
//		img:        rgba,
//		color:      color.Black, // useless
//	}
//	shape := opt.getShape()
//
//	// iterate the matrix to Draw each pixel
//	mat.Iterate(matrix.ROW, func(x int, y int, v matrix.State) {
//		// Draw the block
//		ctx.upperLeft = image.Point{
//			X: x*opt.qrBlockWidth() + _defaultPadding,
//			Y: y*opt.qrBlockWidth() + _defaultPadding,
//		}
//		ctx.lowerRight = image.Point{
//			X: (x+1)*opt.qrBlockWidth() + _defaultPadding,
//			Y: (y+1)*opt.qrBlockWidth() + _defaultPadding,
//		}
//		ctx.color = opt.stateRGBA(v)
//		// DONE(@yeqown): make this abstract to Shapes
//		shape.Draw(ctx)
//	})
//
//	// DONE(@yeqown): add logo image
//	if opt.logoImage() != nil {
//		// Draw logo image into rgba
//		bound := opt.logo.Bounds()
//		upperLeft, lowerRight := bound.Min, bound.Max
//		logoWidth, logoHeight := lowerRight.X-upperLeft.X, lowerRight.Y-upperLeft.Y
//
//		if !validLogoImage(w, h, logoWidth, logoHeight) {
//			log.Printf("w=%d, h=%d, logoW=%d, logoH=%d, logo is over than 1/5 of QRCode \n",
//				w, h, logoWidth, logoHeight)
//			goto done
//		}
//
//		// DONE(@yeqown): calculate the xOffset and yOffset
//		// which point(xOffset, yOffset) should icon upper-left to start
//		xOffset, yOffset := (w-logoWidth)/2, (h-logoHeight)/2
//
//		for posX := upperLeft.X; posX < lowerRight.X; posX++ {
//			for posY := upperLeft.Y; posY < lowerRight.Y; posY++ {
//				rgba.Set(posX+xOffset, posY+yOffset, opt.logo.At(posX, posY))
//			}
//		}
//	}
//done:
//	return rgba
//}

// draw deal QRCode's matrix to be a image.Image
func draw(mat matrix.Matrix, opt *outputImageOptions) image.Image {
	if _debug {
		fmt.Printf("matrix.Width()=%d, matrix.Height()=%d\n", mat.Width(), mat.Height())
	}

	// w as image width, h as image height
	w := mat.Width()*opt.qrBlockWidth() + 2*_defaultPadding
	h := w
	// rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	dc := gg.NewContext(w, h)

	// draw background
	dc.SetColor(opt.backgroundColor())
	dc.DrawRectangle(0, 0, float64(w), float64(h))
	dc.Fill()

	// qrcode block draw context
	ctx := &DrawContext{
		Context:   dc,
		upperLeft: image.Point{},
		w:         opt.qrBlockWidth(),
		h:         opt.qrBlockWidth(),
		color:     color.Black,
	}
	shape := opt.getShape()

	// iterate the matrix to Draw each pixel
	mat.Iterate(matrix.ROW, func(x int, y int, v matrix.State) {
		// Draw the block
		ctx.upperLeft = image.Point{
			X: x*opt.qrBlockWidth() + _defaultPadding,
			Y: y*opt.qrBlockWidth() + _defaultPadding,
		}
		ctx.color = opt.stateRGBA(v)
		// DONE(@yeqown): make this abstract to Shapes
		shape.Draw(ctx)

		//if x == y && _debug {
		//	_ = dc.SavePNG(fmt.Sprintf("./.debug/%d.png", x))
		//}
	})

	//if _debug {
	//	fmt.Printf("save as tmp.png, err=%v\n", dc.SavePNG("./tmp.png"))
	//}

	// DONE(@yeqown): add logo image
	if opt.logoImage() != nil {
		// Draw logo image into rgba
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
		dc.DrawImage(opt.logoImage(), (w-logoWidth)/2, (h-logoHeight)/2)
	}
done:
	return dc.Image()
}

func validLogoImage(qrWidth, qrHeight, logoWidth, logoHeight int) bool {
	return qrWidth >= 5*logoWidth && qrHeight >= 5*logoHeight
}
