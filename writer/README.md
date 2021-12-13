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
- [x] [Terminal output writer](./terminal/README.md)

### How to customize your own writer?

As you can see, the writer is a simple interface, you can implement your own
writer to fit your needs. The prerequisites are:

- Already understand `Writer` interface usage.
- Have understood about paint a picture with go (of course you can refer [standard](./standard) implementation).
- QR Code matrix is a 2D array of `martix.State`, we can simply divide these states into
  binary value (0/1), as the following table:

|  State Expr  | value |                 representation                  |
|:------------:|:-----:|:-----------------------------------------------:|
|  StateFalse  |   0   |              unset (data and etc)               |
|     ZERO     |   0   |              same as `StateFalse`               |
|  StateTrue   |   1   |                   set (data)                    |
|  StateInit   |   1   | not changed since initialized (temporary state) |
| StateVersion |   1   |                set (qr version)                 |
| StateFormat  |   1   |                 set (qr format)                 |
| StateFinder  |   1   |                 set (qr finder)                 |

Now, you can implement your own writer to fit your needs, let's use pseudocode to discuss: 

```text
// define your own writer structure to implement `Writer` interface.
object writer {};

writer.Write(matrix.Matrix):
	// these should be `IMAGE` stream controller, it receives matrix
	// it decide how to print the image. 
	object paint;  
	
	// set use BLACK, unset use WHITE. Or provides matrix.State to colors mapping
	// so that you can control QR Image output intensively. 
	foreach row in matrix:
		foreach column in row:
			if column.State in [StateFalse]:
				// paint a WHITE square block
				paint.draw(x, y, WHITE);
			else:
				// paint a BLACK square block, 
				paint.draw(x, y, BLACK);
	// loop end;
	// output paint;
```