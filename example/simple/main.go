package main

import "github.com/yeqown/go-qrcode"

func main() {
	config := qrcode.Config{
		EncMode: qrcode.EncModeByte,
		EcLevel: qrcode.ErrorCorrectionQuart,
	}
	qrc, err := qrcode.NewWithConfig("github.com/yeqown/qo-qrcode", &config, qrcode.WithQRWidth(40))
	if err != nil {
		panic(err)
	}

	if err = qrc.Save("./a.png"); err != nil {
		panic(err)
	}
}
