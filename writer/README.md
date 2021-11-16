## Writer

Writer folders contains built-in writers those implement `qrcode.Writer`:

```go
// Writer is the interface of a QR code writer, it defines the rule of how to
// `print` the code image from matrix. There's built-in writer to output into
// file, terminal.
type Writer interface {
	// Write writes the code image into itself stream, such as io.Writer,
	// terminal output stream, and etc
	Write(mat matrix.Matrix) error
	
	// Close the writer stream if it exists after QRCode.Save() is called.
	Close() error
}
```

### Implementations

- [x] [Standard output file writer](./standard/README.md)
- [ ] [Terminal output writer](./terminal/README.md)