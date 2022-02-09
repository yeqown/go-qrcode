//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/pkg/errors"
	qrcode "github.com/yeqown/go-qrcode/v2"
	stdw "github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	js.Global().Set("generateQRCode", js.FuncOf(genqrcode))
	fmt.Println("com.github.yeqown.goqrcode.wasm loaded")

	select {}
}

// genqrcode generates a qrcode image and returns the base64 encoded string.
// args should be a string array with length of 2 at most. the first one is the
// content string which will be encoded to qrcode, the second one is the encoding
// option, which is optional.
//
// let result = generateQRCode("content", {
//   qrWidth: 200,
//   qrMargin: 10,
//   qrColor: "#000000",
//   qrBackColor: "#ffffff",
//   encLevel: "H",
//   encVersion: 7,
// })
// more options refer to the `genOption` struct
func genqrcode(_ js.Value, args []js.Value) (v interface{}) {
	var (
		srcContent = ""
		r          = new(genResult)
		opt        = new(genOption)
	)

	defer func() {
		v = r.JSValue()
	}()

	switch len(args) {
	case 0:
		r.setError(errors.New("no args"))
		return
	case 1:
		srcContent = args[0].String()
	case 2:
		fallthrough
	default:
		srcContent = args[0].String()
		opt = optionFromJSValue(args[1])
	}

	//if err := opt.validate(); err != nil {
	//	r.setError(errors.Wrap(err, "invalid option"))
	//	return
	//}

	qrc, err := qrcode.NewWith(srcContent, opt.encodeOptions()...)
	if err != nil {
		r.setError(errors.Wrap(err, "genqrcode qrcode"))
		return
	}

	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w := stdw.NewWithWriter(wr, opt.outputOptions()...)

	err = qrc.Save(w)
	if err != nil {
		r.setError(errors.Wrap(err, "apply output option"))
		return
	}

	// apply image to result
	r.setImage(buf)

	return
}
