package main

import (
	"os"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/file"
)

func main() {
	qrc, err := qrcode.New("with_file_writer")
	if err != nil {
		panic(err)
	}

	w := file.New(os.Stdout)
	if err = qrc.Save(w); err != nil {
		panic(err)
	}
}
