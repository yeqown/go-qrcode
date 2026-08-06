module github.com/yeqown/go-qrcode/cmd/qrcode

go 1.19

require (
	github.com/pkg/errors v0.9.1
	github.com/urfave/cli/v2 v2.27.7
	github.com/yeqown/go-qrcode/v2 v2.3.0
	github.com/yeqown/go-qrcode/writer/standard v1.4.0
	github.com/yeqown/go-qrcode/writer/terminal v1.1.3
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-runewidth v0.0.17 // indirect
	github.com/nsf/termbox-go v1.1.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	golang.org/x/image v0.24.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

//replace (
//    github.com/yeqown/go-qrcode/v2 v2.0.1 => ../../
//    github.com/yeqown/go-qrcode/writer/standard v1.0.0 => ../../writer/standard
//)
