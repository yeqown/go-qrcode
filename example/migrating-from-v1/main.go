package main

import (
	qrv2 "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"

	qrv1 "github.com/yeqown/go-qrcode"
)

func main() {
	// draw a QRCode image with
	// - "https://github.com/yeqown/go-qrcode" source text
	// - circle shape
	// - blue foreground color
	// - white background color
	// - 20 pixel block width
	// - 20 pixel border width
	// - The Highest error correction level

	v1()

	v2()
}

func v1() {
	encodeConfig := &qrv1.Config{
		EcLevel: qrv1.ErrorCorrectionHighest,
		EncMode: qrv1.EncModeAuto,
	}
	qrc, err := qrv1.NewWithConfig("https://github.com/yeqown/go-qrcode", encodeConfig,
		qrv1.WithCircleShape(),
		qrv1.WithFgColorRGBHex("#0000ff"),
		qrv1.WithBgColorRGBHex("#ffffff"),
		qrv1.WithQRWidth(20),
		qrv1.WithBorderWidth(20),
	)
	if err != nil {
		panic(err)
	}

	err = qrc.Save("v1.jpeg")
	if err != nil {
		panic(err)
	}
}

func v2() {
	qrc, err := qrv2.NewWith("https://github.com/yeqown/go-qrcode",
		qrv2.WithErrorCorrectionLevel(qrv2.ErrorCorrectionHighest),
	)
	if err != nil {
		panic(err)
	}

	w, err := standard.New("v2.jpeg",
		standard.WithCircleShape(),
		standard.WithFgColorRGBHex("#0000ff"),
		standard.WithBgColorRGBHex("#ffffff"),
		standard.WithQRWidth(20),
		standard.WithBorderWidth(20),
	)
	if err != nil {
		panic(err)
	}

	err = qrc.Save(w)
	if err != nil {
		panic(err)
	}
}
