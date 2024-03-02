package main

import (
	"github.com/shachardevops/go-qrcode/v2"
	"github.com/shachardevops/go-qrcode/writer/standard"
)

func main() {
	qrc, err := qrcode.New("https://github.com/shachardevops/go-qrcode")
	if err != nil {
		panic(err)
	}

	w0, err := standard.New("./repository_qrcode.png",
		standard.WithHalftone("./test.jpeg"),
		standard.WithQRWidth(21),
	)
	handleErr(err)
	err = qrc.Save(w0)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
