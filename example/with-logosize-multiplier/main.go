package main

import (
	"github.com/yeqown/go-qrcode/writer/standard"

	"github.com/yeqown/go-qrcode/v2"
)

func main() {
	qrc, err := qrcode.NewWith("github.com/yeqown/go-qrcode",
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	)
	if err != nil {
		panic(err)
	}

	// Validates logo size as follows qrWidth >= 2*logoWidth && qrHeight >= 2*logoHeight
	// Instead of default expression qrWidth >= 5*logoWidth && qrHeight >= 5*logoHeight
	w, err := standard.New(
		"./simple.png",
		standard.WithLogoImageFileJPEG("./logo.jpg"),
		standard.WithLogoSizeMultiplier(2),
	)
	if err != nil {
		panic(err)
	}

	if err = qrc.Save(w); err != nil {
		panic(err)
	}
}
