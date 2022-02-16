package halftone

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"

	kit "github.com/yeqown/go-qrcode/image-toolkit"
	"github.com/yeqown/go-qrcode/v2"
)

type option struct {
	mixImage image.Image
	width    int
}

// halftone is an algorithm to generate a halftone image with QR Code image.
// refer to the paper: http://vecg.cs.ucl.ac.uk/Projects/SmartGeometry/halftone_QR/halftoneQR_sigga13.html
type halftone struct {
	option *option
}

func New(inputImg string) *halftone {
	mImg, err := kit.Read(inputImg)
	if err != nil {
		panic(err)
	}

	o := &option{
		mixImage: mImg,
		width:    21, // 3 times
	}

	return &halftone{
		option: o,
	}
}

func (h halftone) Write(mat qrcode.Matrix) error {
	qrImg := draw(mat, h.option)
	return kit.Save(qrImg, "./qr.jpg")
}

func draw(mat qrcode.Matrix, opt *option) image.Image {
	// mask image
	mImg := kit.Binaryzation(
		kit.Scale(opt.mixImage, image.Rect(0, 0, mat.Width()*3, mat.Width()*3), nil),
		60,
	)

	_ = kit.Save(mImg, "./mask.jpg")

	margin := 20
	// closer as image width, h as image height
	w := mat.Width()*opt.width + margin*2
	h := mat.Height()*opt.width + margin*2
	dc := gg.NewContext(w, h)

	// draw background
	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, float64(w), float64(h))
	dc.Fill()

	_ = mImg

	// iterate the matrix to Draw each pixel
	mat.Iterate(qrcode.IterDirection_COLUMN, func(x int, y int, v qrcode.QRValue) {
		dx := x*opt.width + margin
		dy := y*opt.width + margin

		switch v.Type() {
		case qrcode.QRType_DATA:
			if v.IsSet() {
				mixBlock(dc, x, y, margin, opt.width, color.Black, mImg)
				//rect(dc, dx, dy, opt.width, opt.width, color.Black)
			} else {
				mixBlock(dc, x, y, margin, opt.width, color.White, mImg)
				//rect(dc, dx, dy, opt.width, opt.width, color.White)
			}
		default:
			if v.IsSet() {
				rect(dc, dx, dy, opt.width, opt.width, color.Black)
			} else {
				rect(dc, dx, dy, opt.width, opt.width, color.White)
			}
		}
	})

	return dc.Image()
}

func rect(dc *gg.Context, x, y, w, h int, c color.Color) {
	dc.SetColor(c)
	dc.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	dc.Fill()
}

func mixBlock(dc *gg.Context, x, y, margin int, w int, c color.Color, mImg image.Image) {
	var cv color.Color
	dx, dy := x*w+margin, y*w+margin

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cv = mImg.At(x*3+i, y*3+j)
			// cv = color.White
			if i == 1 && j == 1 {
				cv = c
			}
			rect(dc, dx+i*w/3, dy+j*w/3, w/3, w/3, cv)
		}
	}
}

func (h halftone) Close() error {
	return nil
}
