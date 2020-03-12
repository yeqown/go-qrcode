package main

import (
	"flag"
	"fmt"

	qrcode "github.com/yeqown/go-qrcode"
)

var (
	contentFlag = flag.String("c", "https://github.com/yeqown/go-qrcode", "input your content to encode")
	outputFlag  = flag.String("o", "./qrcode.jpeg", "output filename")
)

func main() {
	flag.Parse()

	qrc, err := qrcode.New(*contentFlag)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	// save file
	if err := qrc.Save(*outputFlag); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
