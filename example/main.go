package main

import (
	"fmt"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	repo()

}

func repo() {
	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	w, err := standard.New("../assets/repo-qrcode.jpeg")
	if err != nil {
		fmt.Printf("standard.New failed: %v", err)
		return
	}

	// save file
	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
