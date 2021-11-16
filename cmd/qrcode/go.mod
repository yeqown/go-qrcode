module qrcode

go 1.17


require (
	github.com/yeqown/go-qrcode/writer/standard v1.0.0
	github.com/yeqown/go-qrcode/v2 v2.0.0
)

replace (
    github.com/yeqown/go-qrcode/v2 v2.0.0 => ../../
    github.com/yeqown/go-qrcode/writer/standard v1.0.0 => ../../writer/standard
)