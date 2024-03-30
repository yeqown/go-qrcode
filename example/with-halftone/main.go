package main

import (
	"flag"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

var (
	transparent = flag.Bool("transparent", false, "set background to transparent")
)

func main() {
	flag.Parse()

	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		panic(err)
	}

	options := []standard.ImageOption{
		standard.WithHalftone("./test.jpeg"),
		standard.WithQRWidth(21),
	}
	filename := "./halftone-qr.png"

	if *transparent {
		options = append(
			options,
			standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
			standard.WithBgTransparent(),
		)
		filename = "./halftone-qr-transparent.png"
	}

	w0, err := standard.New(filename, options...)
	handleErr(err)
	err = qrc.Save(w0)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
