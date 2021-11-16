package main

import (
	"flag"
	"fmt"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
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

	w, err := standard.New(*outputFlag)
	if err != nil {
		fmt.Printf("standard writer failed: %v", err)
		return
	}

	// save file
	if err := qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
