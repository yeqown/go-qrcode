## WebAssembly support

This part provide support to compile and run WebAssembly code.

```javascript
let option = {
	encodeVersion: 0,   // 0 - 40
	encodeMode: 2,      // 0 - 3
	encodeECLevel: "Q", // L, M, Q, H
	
	outputBgColor: "#123123",   // #000000 - #ffffff
	outputBgTransparent: false, // true - false
	outputQrColor: "#666666",   // #000000 - #ffffff
	outputQrWidth: 20,          // 0 - 256
	outputCircleShape: true,    // true - false
	outputImageEncoder: "png",  // png, jpeg, jpg
	outputMargin: 20,           // 0 - 256
}
let r = generateQRCode("content", option)
// output:
// {
    "success": true,
    "error": "",
    "base64EncodedImage": "iVBORw0KGgoAAAANSUhEUgAAAmwAAAJ... more"
}
```

## build wasm binary


```bash
cd $PATH/to/go-qrcode/cmd/wasm
GOOS=js GOARCH=wasm go build -o com.github.yeqown.goqrcode.wasm .
```