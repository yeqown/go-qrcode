## WebAssembly Example

This example demonstrates how to use the `go-qrcode/wasm` pre-compiled `wasm` to generate
QRCode image in Web browser.

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
cp "$PATH/go-qrcode/wasm/com.github.yeqown.goqrcode.wasm" .
# then serve the index.html in browser
python3 -m http.server 
# it serves the index.html on http://localhost:8000, you can use another way too.
```

after that, you can visit http://localhost:8000 to see the demo.