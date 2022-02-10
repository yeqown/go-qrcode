//go:build js && wasm

package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"syscall/js"

	"github.com/yeqown/go-qrcode/v2"
	stdw "github.com/yeqown/go-qrcode/writer/standard"
)

// encodeOption refers to github.com/yeqown/go-qrcode/v2.encodingOption
type encodeOption struct {
	version int    // encodeVersion
	mode    uint8  // encodeMode specifies which encMode to use (0/1/2/3).
	ecLevel string // encodeECLevel specifies which ecLevel to use (L/M/Q/H)
}

// outputOption refers to github.com/yeqown/go-qrcode/writer/standard.outputImageOptions
type outputOption struct {
	bgColor       string // outputBgColor is the background color of the QR code image.
	bgTransparent bool   // outputBgTransparent indicates whether the background color is transparent.
	qrColor       string // outputQrColor is the foreground color of the QR code.
	qrWidth       uint8  // outputQrWidth is the width of the QR code.
	circleShape   bool   // outputCircleShape indicates whether to draw the qr block in circle shape.
	imageEncoder  string // outputImageEncoder specifies file format would be encoded the QR image. (jpg/jpeg/png)
	margin        int    // outputMargin is the border width of the output image.
}

// genOption is a type of option for generating code.
type genOption struct {
	encodeOption
	outputOption
}

// optionFromJSValue converts js.Value to genOption.
func optionFromJSValue(option js.Value) *genOption {
	//fmt.Println(option.Type().String())
	if option.IsNull() || option.IsUndefined() || option.Type().String() != "object" {
		// option must be an object, otherwise return default empty option.
		return nil
	}

	return &genOption{
		encodeOption: encodeOption{
			version: option.Get("encodeVersion").Int(),
			mode:    uint8(option.Get("encodeMode").Int()),
			ecLevel: option.Get("encodeECLevel").String(),
		},
		outputOption: outputOption{
			bgColor:       option.Get("outputBgColor").String(),
			bgTransparent: option.Get("outputBgTransparent").Bool(),
			qrColor:       option.Get("outputQrColor").String(),
			qrWidth:       uint8(option.Get("outputQrWidth").Int()),
			circleShape:   option.Get("outputCircleShape").Bool(),
			imageEncoder:  option.Get("outputImageEncoder").String(),
			margin:        option.Get("outputMargin").Int(),
		},
	}
}

// validate genOption.
func (o *genOption) validate() error {
	if o == nil {
		return nil
	}

	return nil
}

func (o *genOption) encodeOptions() []qrcode.EncodeOption {
	if o == nil {
		return nil
	}

	out := make([]qrcode.EncodeOption, 0, 4)

	if o.encodeOption.version != 0 {
		out = append(out, qrcode.WithVersion(o.encodeOption.version))
	}

	switch o.encodeOption.mode {
	case uint8(qrcode.EncModeAlphanumeric):
		out = append(out, qrcode.WithEncodingMode(qrcode.EncModeAlphanumeric))
	case uint8(qrcode.EncModeNumeric):
		out = append(out, qrcode.WithEncodingMode(qrcode.EncModeNumeric))
	case uint8(qrcode.EncModeByte):
		out = append(out, qrcode.WithEncodingMode(qrcode.EncModeByte))
	case uint8(qrcode.EncModeJP):
		out = append(out, qrcode.WithEncodingMode(qrcode.EncModeJP))
	}

	switch o.encodeOption.ecLevel {
	case "L":
		out = append(out, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow))
	case "M":
		out = append(out, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	case "Q":
		out = append(out, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart))
	case "H":
		out = append(out, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest))
	}

	return out
}

func (o *genOption) outputOptions() []stdw.ImageOption {
	if o == nil {
		return nil
	}

	out := make([]stdw.ImageOption, 0, 8)

	if o.outputOption.bgColor != "" {
		out = append(out, stdw.WithBgColorRGBHex(o.outputOption.bgColor))
	}
	if o.outputOption.bgTransparent {
		out = append(out, stdw.WithBgTransparent())
	}
	if o.outputOption.qrColor != "" {
		out = append(out, stdw.WithFgColorRGBHex(o.outputOption.qrColor))
	}
	if o.outputOption.qrWidth != 0 {
		out = append(out, stdw.WithQRWidth(o.outputOption.qrWidth))
	}
	if o.outputOption.circleShape {
		out = append(out, stdw.WithCircleShape())
	}

	switch o.outputOption.imageEncoder {
	case "jpg", "jpeg":
		out = append(out, stdw.WithBuiltinImageEncoder(stdw.JPEG_FORMAT))
	case "png":
		fallthrough
	default:
		out = append(out, stdw.WithBuiltinImageEncoder(stdw.PNG_FORMAT))
	}

	if o.outputOption.margin != 0 {
		out = append(out, stdw.WithBorderWidth(o.outputOption.margin))
	}

	return out
}

// genResult is result contains generated code image or error message.
type genResult struct {
	Success            bool   `json:"success"`
	Error              string `json:"error"`
	Base64EncodedImage string `json:"base64EncodedImage"`
}

func (r *genResult) setError(err error) {
	r.Success = false
	r.Error = err.Error()
}

func (r *genResult) setImage(buf *bytes.Buffer) {
	r.Success = true
	r.Base64EncodedImage = base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (r *genResult) JSValue() js.Value {
	return js.ValueOf(map[string]any{
		"success":            r.Success,
		"error":              r.Error,
		"base64EncodedImage": r.Base64EncodedImage,
	})
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
