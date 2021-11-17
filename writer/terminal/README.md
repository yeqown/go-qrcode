## Terminal 

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/yeqown/go-qrcode/writer/standard)

Standard Writer is a writer that is used to draw QR Code image into terminal.

### Usage

```go
package main

import (
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/terminal"
)

func main() {
	qrc, _ := qrcode.New("withTerminalWriter")

	w := terminal.New()
	
	if err := qrc.Save(w); err != nil {
		panic(err)
	}
}
```

### Option

> 🤪 Do not support any option yet.