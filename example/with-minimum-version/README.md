# QR Code with Minimum Version Example

This example demonstrates how to use the `WithMinimumVersion` option to ensure that a QR code has at least a specific version level.

## What is QR Code Version?

QR code versions range from 1 to 40, with higher versions supporting more data capacity and having larger dimensions:
- Version 1: 21x21 modules
- Version 9: 53x53 modules  
- Version 40: 177x177 modules

The formula is: dimension = (version × 4) + 17

## Why Use Minimum Version?

Sometimes you want to ensure a QR code has a minimum size or capacity, even if the content would fit in a smaller version:

1. **Consistency**: Keep QR codes the same size across different content lengths
2. **Future-proofing**: Reserve capacity for potential content expansion
3. **Readability**: Larger QR codes are easier to scan from a distance
4. **Design requirements**: Match specific design or layout requirements

## Usage

```go
qrc, err := qrcode.NewWith("your-text-here",
    qrcode.WithMinimumVersion(9),
    qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
)
```

## Running the Example

```bash
go run main.go
```

This will generate two QR codes:
1. `qrcode_minimum_version_9.png` - A QR code with minimum version 9 (53x53 modules)
2. `qrcode_auto_version.png` - The same content with automatic version selection (much smaller)

## Note

If the content naturally requires a version higher than the minimum, the higher version will be used. The minimum version only applies when the automatic version would be lower.
