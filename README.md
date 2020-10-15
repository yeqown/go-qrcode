# go-qrcode #

<img src="./qrcode.jpeg" width="100px" align="right"/>
QR code (abbreviated from Quick Response Code) is the trademark for a type of matrix barcode (or two-dimensional barcode) first designed in 1994 for the automotive industry in Japan. A barcode is a machine-readable optical label that contains information about the item to which it is attached. A QR code uses four standardized encoding modes (numeric, alphanumeric, byte/binary, and kanji) to store data efficiently; extensions may also be used

### Install

```sh
go get -u github.com/yeqown/go-qrcode
```

### Usage

link to [CODE](./example/main.go)
```go
package main

import (
	"fmt"

	qrcode "github.com/yeqown/go-qrcode"
)

func main() {
	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}

	// save file
	if err := qrc.Save("../testdata/repo-qrcode.jpeg"); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
```

### Documention

Jump to [godoc.org/github/yeqown/go-qrcode](https://godoc.org/github.com/yeqown/go-qrcode)

### Links

* [QR Code tutori](https://www.thonky.com/qr-code-tutorial/)
* [QRCode Wiki](https://en.wikipedia.org/wiki/QR_code)
* [二维码详解（QR Code）](https://zhuanlan.zhihu.com/p/21463650)
* [数据编码](https://zhuanlan.zhihu.com/p/25432676)