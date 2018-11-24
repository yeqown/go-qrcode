package main

import (
	"fmt"

	qrcode "github.com/yeqown/go-qrcode"
)

func main() {
	// 配置文件，默认在repo的根路径下
	qrcode.SetVersionCfgFile("../versionCfg.json")

	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}

	// save file
	if err := qrc.Save("../testdata/repo-qrcode.jpeg"); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
