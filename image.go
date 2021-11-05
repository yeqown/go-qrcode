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

	defer func(f *os.File) {
		err = f.Close()
	}(f)

	return drawTo(f, m, opt)
}

// drawTo save image into io.Writer
func drawTo(w io.Writer, m matrix.Matrix, opt *outputImageOptions) (err error) {
	if opt == nil {
		opt = defaultOutputImageOption()
	}

	img := draw(m, opt)

	// DONE(@yeqown): support file format specified config option
	if err = opt.imageEncoder.Encode(w, img); err != nil {
		err = fmt.Errorf("jpeg.Encode got err: %v", err)
	}

	return
}

// draw deal QRCode's matrix to be an image.Image. Notice that if anyone changed this function,
// please also check the function outputImageOptions.preCalculateAttribute().
func draw(mat matrix.Matrix, opt *outputImageOptions) image.Image {
	if _debug {
		fmt.Printf("matrix.Width()=%d, matrix.Height()=%d\n", mat.Width(), mat.Height())
	}

	top, right, bottom, left := opt.borderWidths[0], opt.borderWidths[1], opt.borderWidths[2], opt.borderWidths[3]
	// w as image width, h as image height
	w := mat.Width()*opt.qrBlockWidth() + left + right
	h := mat.Height()*opt.qrBlockWidth() + top + bottom
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
			X: x*opt.qrBlockWidth() + left,
			Y: y*opt.qrBlockWidth() + top,
		}
		ctx.color = opt.stateRGBA(v)
		// DONE(@yeqown): make this abstract to Shapes

		switch v {
		case matrix.StateFinder:
			shape.DrawFinder(ctx)
		default:
			shape.Draw(ctx)
		}

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

// Attribute contains basic information of generated image.
type Attribute struct {
	// width and height of image
	W, H int
	// in the order of "top, right, bottom, left"
	Borders [4]int
	// the length of  block edges
	BlockWidth int
}
