package main

import (
	"fmt"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	// Example 1: Generate a QR code with minimum version 9
	// Even though the text "Hello" would fit in a smaller version,
	// the QR code will be at least version 9
	qrc, err := qrcode.NewWith("Hello",
		qrcode.WithMinimumVersion(9),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	)
	if err != nil {
		panic(fmt.Sprintf("could not generate QRCode: %v", err))
	}

	w, err := standard.New("./qrcode_minimum_version_9.png", standard.WithQRWidth(20))
	if err != nil {
		panic(fmt.Sprintf("could not create writer: %v", err))
	}

	if err = qrc.Save(w); err != nil {
		panic(fmt.Sprintf("could not save QRCode: %v", err))
	}

	fmt.Printf("QR code generated with minimum version 9\n")
	fmt.Printf("Dimension: %d x %d\n", qrc.Dimension(), qrc.Dimension())
	fmt.Printf("Expected dimension for version 9: %d x %d\n", 9*4+17, 9*4+17)

	// Example 2: Compare with automatic version selection
	qrc2, err := qrcode.New("Hello")
	if err != nil {
		panic(fmt.Sprintf("could not generate QRCode: %v", err))
	}

	w2, err := standard.New("./qrcode_auto_version.png", standard.WithQRWidth(20))
	if err != nil {
		panic(fmt.Sprintf("could not create writer: %v", err))
	}

	if err = qrc2.Save(w2); err != nil {
		panic(fmt.Sprintf("could not save QRCode: %v", err))
	}

	fmt.Printf("\nFor comparison, QR code with automatic version:\n")
	fmt.Printf("Dimension: %d x %d (much smaller)\n", qrc2.Dimension(), qrc2.Dimension())
}
