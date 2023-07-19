module example

go 1.17

require (
	github.com/yeqown/go-qrcode v1.5.10
	github.com/yeqown/go-qrcode/v2 v2.2.2
	github.com/yeqown/go-qrcode/writer/standard v1.1.1
	github.com/yeqown/go-qrcode/writer/terminal v1.0.0-beta
)

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/nsf/termbox-go v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	golang.org/x/image v0.5.0 // indirect
)

replace (
	github.com/yeqown/go-qrcode/v2 => ../
	github.com/yeqown/go-qrcode/writer/standard => ../writer/standard
	github.com/yeqown/go-qrcode/writer/terminal v1.0.0-beta => ../writer/terminal
)
