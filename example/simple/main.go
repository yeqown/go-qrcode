package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	qrc, err := qrcode.NewWith("github.com/yeqown/go-qrcode",
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	)
	if err != nil {
		panic(err)
	}

	// save into file
	w, err := standard.New("./simple.png", standard.WithQRWidth(40))
	if err != nil {
		panic(err)
	}
	if err = qrc.Save(w); err != nil {
		panic(err)
	}

	// get bytes
	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w2 := standard.NewWithWriter(wr, standard.WithQRWidth(40))
	if err = qrc.Save(w2); err != nil {
		panic(err)
	}

	// copy base64 content to converter, then you can get the picture
	// https://codebeautify.org/base64-to-image-converter
	fmt.Printf("encoded base64 image:\n%s", base64.StdEncoding.EncodeToString(buf.Bytes()))
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
