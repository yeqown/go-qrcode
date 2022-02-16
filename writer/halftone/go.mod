module github.com/yeqown/go-qrcode/writer/halftone

go 1.18

require (
	github.com/fogleman/gg v1.3.0
	github.com/yeqown/go-qrcode/image-toolkit v1.0.0
	github.com/yeqown/go-qrcode/v2 v2.1.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
)

replace (
	github.com/yeqown/go-qrcode/image-toolkit v1.0.0 => ../../image-toolkit
	github.com/yeqown/go-qrcode/v2 => ../../
)
