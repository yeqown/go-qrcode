module github.com/yeqown/go-qrcode/cmd/qrcode

go 1.17

require (
	github.com/pkg/errors v0.9.1
	github.com/urfave/cli/v2 v2.3.0
	github.com/yeqown/go-qrcode/v2 v2.2.2
	github.com/yeqown/go-qrcode/writer/standard v1.2.0
	github.com/yeqown/go-qrcode/writer/terminal v1.1.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/nsf/termbox-go v1.1.1 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	golang.org/x/image v0.5.0 // indirect
)

//replace (
//    github.com/yeqown/go-qrcode/v2 v2.0.1 => ../../
//    github.com/yeqown/go-qrcode/writer/standard v1.0.0 => ../../writer/standard
//)
