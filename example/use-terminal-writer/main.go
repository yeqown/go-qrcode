package main

import (
	"github.com/shachardevops/go-qrcode/v2"
	"github.com/shachardevops/go-qrcode/writer/terminal"
)

func main() {
	qrc, err := qrcode.New("with_terminal")
	if err != nil {
		panic(err)
	}

	w := terminal.New()
	if err = qrc.Save(w); err != nil {
		panic(err)
	}

	//time.Sleep(5 * time.Second)
}
