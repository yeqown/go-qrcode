package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"github.com/yeqown/go-qrcode/writer/terminal"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var copyright = `Copyright (c) 2018 yeqown

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
`

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "qrcode"
	app.Description = "QR code generator"
	app.Version = "v2.0.1"
	app.Authors = []*cli.Author{
		{
			Name:  "yeqown",
			Email: "yeqown@gmail.com",
		},
	}
	app.Usage = "qrcode [options] [source text]"
	app.Copyright = copyright
	app.Flags = prepareFlags()
	app.Action = func(c *cli.Context) error {
		genCtx := parseGenerateContextFrom(c)
		return generate(genCtx)
	}

	return app
}

func generate(ctx *generateContext) error {
	qrc, err := qrcode.New(ctx.text)
	if err != nil {
		return errors.Wrap(err, "generate QRCode failed")
	}

	var w qrcode.Writer

	// construct a writer based generateContext
	switch ctx.mode {
	case writerMode_TERMINAL:
		w = terminal.New()
	default:
		w, err = standard.New(ctx.FOO.output, ctx.FOO.applyOptions()...)
	}
	if err != nil {
		return errors.Wrap(err, "initialize writer failed")
	}

	return qrc.Save(w)
}

type writerMode uint8

const (
	_ writerMode = iota
	writerMode_FILE
	writerMode_TERMINAL
)

// generateContext generate qrcode from context
type generateContext struct {
	text string
	mode writerMode
	FOO  *fileOutputOptions
	TOO  *terminalOutputOptions
}

type fileOutputOptions struct {
	output        string
	outputSuffix  string
	blockSize     uint8
	borders       [4]int
	isCircleShape bool
	transparent   bool
	halftoneImage string
}

func (foo fileOutputOptions) applyOptions() []standard.ImageOption {
	options := []standard.ImageOption{
		standard.WithQRWidth(foo.blockSize),
		standard.WithBorderWidth(foo.borders[:]...),
	}

	switch foo.outputSuffix {
	case "png":
		options = append(options, standard.WithBuiltinImageEncoder(standard.PNG_FORMAT))
	case "jpg", "jpeg":
		fallthrough
	default:
		options = append(options, standard.WithBuiltinImageEncoder(standard.JPEG_FORMAT))
	}

	if foo.isCircleShape {
		options = append(options, standard.WithCircleShape())
	}

	if foo.transparent {
		options = append(options, standard.WithBgTransparent())
	}

	if foo.halftoneImage != "" {
		options = append(options, standard.WithHalftone(foo.halftoneImage))
	}

	return options
}

type terminalOutputOptions struct{}

func parseGenerateContextFrom(c *cli.Context) *generateContext {
	genCtx := &generateContext{
		text: c.Args().First(),
		mode: writerMode_FILE,
		FOO: &fileOutputOptions{
			output:        c.String("output"),
			outputSuffix:  strings.TrimPrefix(filepath.Ext(c.String("output")), "."),
			blockSize:     uint8(c.Uint("block")),
			borders:       [4]int{},
			isCircleShape: c.Bool("circle"),
			transparent:   c.Bool("transparent"),
			halftoneImage: c.String("halftone"),
		},
		TOO: &terminalOutputOptions{},
	}

	// writer mode
	if c.Bool("terminal") {
		genCtx.mode = writerMode_TERMINAL
	}

	// parse borders
	borders := c.String("borders")
	arr := strings.Split(borders, ",")
	if len(arr) != 4 {
		panic("invalid borders format, want: uint8,uint8,uint8,uint8")
	}
	for i, s := range arr {
		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			panic("invalid borders format, want: uint8")
		}
		genCtx.FOO.borders[i] = int(v)
	}

	return genCtx
}

func prepareFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "terminal",
			Usage:       "--terminal",
			Value:       false,
			DefaultText: "false",
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "--output=<output file>",
			Value:       "qrcode.jpg",
			DefaultText: "qrcode.jpg",
		},
		&cli.UintFlag{
			Name:        "block",
			Aliases:     []string{"s"},
			Usage:       "--block=<block size>",
			Value:       5,
			DefaultText: "5",
		},
		&cli.StringFlag{
			Name:        "borders",
			Aliases:     []string{"b"},
			Usage:       "--borders=<borders>",
			Value:       "0,0,0,0",
			DefaultText: "0,0,0,0",
		},
		&cli.BoolFlag{
			Name:        "circle",
			Usage:       "--circle",
			Value:       false,
			DefaultText: "false",
		},
		&cli.BoolFlag{
			Name:        "transparent",
			Usage:       "--transparent",
			Value:       false,
			DefaultText: "false",
		},
		&cli.StringFlag{
			Name:        "halftone",
			Usage:       "--halftone=<halftone image>",
			Value:       "",
			DefaultText: "",
		},
	}
}
