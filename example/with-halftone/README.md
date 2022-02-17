## with-halftone

> 🚧 Notice that could not use another shape except for the rectangle since
> the gap between qr blocks could mislead the recognizer.

```go
func main() {
	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
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
```

#### input image

<img src="./test.jpeg">

#### output image

<img src="./halftone-qr.jpeg">