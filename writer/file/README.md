## File

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/yeqown/go-qrcode/writer/file)

File Writer is a writer used to draw QR Code images into files using the characters ▀, ▄, █, and space.

### Usage

```go
package main

import (
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/file"
)

func main() {
	qrc, _ := qrcode.New("with_file_writer")

	w := file.New(os.Stdout)

	if err := qrc.Save(w); err != nil {
		panic(err)
	}
}
```