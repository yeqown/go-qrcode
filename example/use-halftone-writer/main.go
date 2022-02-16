package main

import (
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/halftone"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		panic(err)
	}

	w0, err := standard.New("./standard.jpeg")
	handleErr(err)
	err = qrc.Save(w0)
	handleErr(err)

	w := halftone.New("./test.jpeg")
	err = qrc.Save(w)
	handleErr(err)

	//time.Sleep(5 * time.Second)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
