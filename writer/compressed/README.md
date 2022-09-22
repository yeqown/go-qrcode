## Compressed Writer

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/yeqown/go-qrcode/writer/compressed)

Compressed Writer is a writer that is used to draw QR Code image on a very small scale.
Since the compressed writer would only use a two-tone palette to generate the image,
and automatically compressed in PNG format.

Check codes for more details about compress principle.

### Usage

```go
option := compressed.Option{
	Padding:   4, // padding pixels around the qr code.
	BlockSize: 1, // block pixels which represents a bit data.
}

w, err := compressed.New(name, &option)
	if err != nil {
	panic(err)
}

if err := qrc.Save(w); err != nil {
	panic(err)
}
```