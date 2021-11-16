package standard

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/v2/matrix"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

var _ qrcode.Writer = (*Writer)(nil)

var (
	ErrNilWriter = errors.New("nil writer")
)

// Writer is a writer that writes QR Code to io.Writer.
type Writer struct {
	option *outputImageOptions

	closer io.WriteCloser
}

// New creates a standard writer.
func New(filename string, opts ...ImageOption) (*Writer, error) {
	if _, err := os.Stat(filename); err != nil && os.IsExist(err) {
		// custom path got: "file exists"
		log.Printf("could not find path: %s, then save to %s", filename, _defaultFilename)
		filename = _defaultFilename
	}

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "create file failed")
	}

	return NewWithWriter(fd, opts...), nil
}

func NewWithWriter(writeCloser io.WriteCloser, opts ...ImageOption) *Writer {
	dst := defaultOutputImageOption()
	for _, opt := range opts {
		opt.apply(dst)
	}

	if writeCloser == nil {
		panic("writeCloser could not be nil")
	}

	return &Writer{
		option: dst,
		closer: writeCloser,
	}
}

const (
	_defaultFilename = "default.jpeg"
	_defaultPadding  = 40
)

func (w Writer) Write(mat matrix.Matrix) error {
	defer func() {
		_ = w.Close()
	}()

	return drawTo(w.closer, mat, w.option)
}

func (w Writer) Close() error {
	if w.closer == nil {
		return nil
	}

	if err := w.closer.Close(); !errors.Is(err, os.ErrClosed) {
		return err
	}

	return nil
}

func (w Writer) Attribute(dimension int) *Attribute {
	return w.option.preCalculateAttribute(dimension)
}

func drawTo(w io.Writer, mat matrix.Matrix, option *outputImageOptions) (err error) {
	if option == nil {
		option = defaultOutputImageOption()
	}

	if w == nil {
		return ErrNilWriter
	}

	img := draw(mat, option)

	// DONE(@yeqown): support file format specified config option
	if err = option.imageEncoder.Encode(w, img); err != nil {
		err = fmt.Errorf("imageEncoder.Encode failed: %v", err)
	}

	return
}

// draw deal QRCode's matrix to be an image.Image. Notice that if anyone changed this function,
// please also check the function outputImageOptions.preCalculateAttribute().
func draw(mat matrix.Matrix, opt *outputImageOptions) image.Image {
	//if v2._debug {
	//	fmt.Printf("matrix.Width()=%d, matrix.Height()=%d\n", mat.Width(), mat.Height())
	//}

	top, right, bottom, left := opt.borderWidths[0], opt.borderWidths[1], opt.borderWidths[2], opt.borderWidths[3]
	// closer as image width, h as image height
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
			log.Printf("closer=%d, h=%d, logoW=%d, logoH=%d, logo is over than 1/5 of QRCode \n",
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
