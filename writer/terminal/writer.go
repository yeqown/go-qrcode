package terminal

import (
	"github.com/shachardevops/go-qrcode/v2"

	"github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

var _ qrcode.Writer = (*Writer)(nil)

// Writer implements qrcode.Writer based on termbox to print QRCode into
// terminal / console.
type Writer struct{}

func New() *Writer {
	w := &Writer{}
	w.init()

	return w
}

func (w Writer) init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
}

func (w Writer) preDraw(width, height, padding int, bg termbox.Attribute) {
	for i := 0; i < width+2*padding; i++ {
		for j := 0; j < height+2*padding; j++ {
			w.drawBlock(i, j, 0, bg, bg)
		}
	}
}

// drawBlock draws a block at (x, y) with fg and bg colors.
// each block takes 2 times width of one character terminal, it looks like: ██
func (w Writer) drawBlock(x, y, padding int, fg termbox.Attribute, bg termbox.Attribute) {
	x1, y1 := x*2+2*padding, y+padding
	x2, y2 := x1+1, y1

	termbox.SetCell(x1, y1, '█', fg, bg)
	termbox.SetCell(x2, y2, '█', fg, bg)
}

func (w Writer) Write(mat qrcode.Matrix) error {
	//width, height, whratio := terminalSize()
	//_ = width
	//_ = height
	//_ = whratio

	ww, hh := mat.Width(), mat.Height()
	bg := termbox.ColorWhite
	fg := termbox.ColorBlack

	padding, curRow := 1, 0
	w.preDraw(ww, hh, padding, bg)
	mat.Iterate(qrcode.IterDirection_ROW, func(x int, y int, state qrcode.QRValue) {
		if state.IsSet() {
			fg = termbox.ColorBlack
		} else {
			fg = termbox.ColorWhite
		}

		w.drawBlock(x, y, padding, fg, bg)
		curRow = y
	})

	printTip(curRow + 2*padding + 1 + 1)
	return hold()
}

func printTip(y int) {
	tip := "Press any key to quit."
	x := 0
	for _, r := range tip {
		w := runewidth.RuneWidth(r)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
			w = 1
		}
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
		x += w
	}
}

func hold() error {
	if err := termbox.Flush(); err != nil {
		return err
	}

wait:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break wait
		}
	}

	return nil
}

func (w Writer) Close() error {
	termbox.Close()

	return nil
}
