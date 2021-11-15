package main

import (
	"github.com/yeqown/go-qrcode/writer/standard"

	"github.com/yeqown/go-qrcode/v2"
)

func main() {
	qrc, err := qrcode.NewWith("github.com/yeqown/qo-qrcode",
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	)
	if err != nil {
		panic(err)
	}

	w, err := standard.New("./simple.png", standard.WithQRWidth(40))
	if err != nil {
		panic(err)
	}

	if err = qrc.Save(w); err != nil {
		panic(err)
	}
}
