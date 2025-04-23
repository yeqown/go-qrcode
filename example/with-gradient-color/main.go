package main

import (
	"image/color"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	qrc, err := qrcode.NewWith("https://github.com/yeqown/go-qrcode", qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow))
	if err != nil {
		panic(err)
	}
	g := standard.NewGradient(45, []standard.ColorStop{
		{
			T:     0,
			Color: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			T:     0.5,
			Color: color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			T:     1,
			Color: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		},
	}...)
	w, err := standard.New("./with-gradient-color/smaller.png",
		standard.WithFgGradient(g),
	)
	if err != nil {
		panic(err)
	}

	err = qrc.Save(w)
	if err != nil {
		panic(err)
	}
}
