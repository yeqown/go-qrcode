package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// createQRWithQRWidth generates a QR code using the WithQRWidth option.
func createQRWithQRWidth(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	options := []standard.ImageOption{
		standard.WithQRWidth(8),
	}
	writer, err := standard.New("../assets/example/qrcode_with_qrwidth.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

// createQRWithCircleShape generates a QR code using the WithCircleShape option.
func createQRWithCircleShape(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	options := []standard.ImageOption{
		standard.WithCircleShape(),
	}
	writer, err := standard.New("../assets/example/qrcode_with_circleshape.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

// createQRWithBorderWidth generates a QR code using the WithBorderWidth option.
func createQRWithBorderWidth(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	options := []standard.ImageOption{
		standard.WithBorderWidth(10),
	}
	writer, err := standard.New("../assets/example/qrcode_with_borderwidth.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

// createQRWithBgTransparent generates a QR code using the WithBgTransparent option.
func createQRWithBgTransparent(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	options := []standard.ImageOption{
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithBgTransparent(),
	}
	writer, err := standard.New("../assets/example/qrcode_with_bgtransparent.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

// createQRWithBgColor generates a QR code using the WithBgColor option.
func createQRWithBgFgColor(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	options := []standard.ImageOption{
		standard.WithFgColor(color.RGBA{135, 206, 235, 255}), // standard.WithBgColorRGBHex(),
		standard.WithBgColor(color.RGBA{124, 252, 0, 255}),   // standard.WithFgColorRGBHex(),
	}
	writer, err := standard.New("../assets/example/qrcode_with_bgfgcolor.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

// createQRWithFgGradient generates a QR code using the WithFgGradient option.
func createQRWithFgGradient(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	stops := []standard.ColorStop{
		{Color: color.RGBA{255, 0, 0, 255}, T: 0.0},
		{Color: color.RGBA{0, 255, 0, 255}, T: 0.5},
		{Color: color.RGBA{0, 0, 255, 255}, T: 1},
	}
	// Create a linear gradient
	gradient := standard.NewGradient(45, stops...)

	options := []standard.ImageOption{
		standard.WithFgGradient(gradient),
	}
	writer, err := standard.New("../assets/example/qrcode_with_fggradient.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

// createQRWithHalftone generates a QR code using the WithHalftone option.
func createQRWithHalftone(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	// Please replace with the actual path to the halftone image.
	halftonePath := "../assets/example/monna-lisa.png"
	if _, err := os.Stat(halftonePath); os.IsNotExist(err) {
		fmt.Printf("halftone image file %s not found\n", halftonePath)
		return
	}

	options := []standard.ImageOption{
		standard.WithHalftone(halftonePath),
		standard.WithQRWidth(21),
	}
	writer, err := standard.New("../assets/example/qrcode_with_halftone.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

// createQRWithLogo generates a QR code using the WithLogo option.
func createQRWithLogo(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	options := []standard.ImageOption{
		standard.WithLogoImageFileJPEG("../assets/example/logo.jpeg"),
	}
	writer, err := standard.New("../assets/example/qrcode_with_logo.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}

	defer writer.Close()
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}

func main() {
	content := "https://github.com/yeqown/go-qrcode"

	createQRWithQRWidth(content)
	createQRWithCircleShape(content)
	createQRWithBorderWidth(content)
	createQRWithBgTransparent(content)
	createQRWithBgFgColor(content)
	createQRWithFgGradient(content)
	createQRWithHalftone(content)
	createQRWithLogo(content)

	fmt.Println("All QR codes generated successfully.")
}
