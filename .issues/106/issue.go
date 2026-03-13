/*
 * Link:	https://github.com/yeqown/go-qrcode/issues/106
 * Title:	Feature: Add Kanji encoding mode support
 * Author:	fdelbos(https://github.com/fdelbos)
 */

package main

import (
	"fmt"

	yeqown "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/compressed"
)

/*
See https://github.com/yeqown/go-qrcode/issues/106
Feature: Add Kanji encoding mode support for QR codes

Results:
Content length:     3   // source text length (3 Kanji characters)
qr-kanji.png:       <size> bytes
*/
func main() {
	// Pure Kanji text with explicit Kanji mode
	// But Shift-JIS only supports Kanji characters, not full-width alphanumeric,
	// so we can't encode `https://google.com` in Kanji mode
	content := "日本語"

	fmt.Printf("Content length: %d\n", len([]rune(content)))

	qrc, err := yeqown.NewWith(content,
		yeqown.WithEncodingMode(yeqown.EncModeKanji),
	)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	// Save to file
	option := &compressed.Option{
		Padding:   4,
		BlockSize: 1,
	}
	w, err := compressed.New("qr-kanji.png", option)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if err = qrc.Save(w); err != nil {
		panic(err)
	}

	fmt.Println("QR code saved to qr-kanji.png")

	// Note: If your input might contain non-Kanji characters, use EncModeAuto:
	// qrc, err := yeqown.NewWith(anyText,
	// 	yeqown.WithEncodingMode(yeqown.EncModeAuto),
	// )
}
