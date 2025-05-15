package shapes

import "github.com/yeqown/go-qrcode/writer/standard"

// RoundedFinder returns a function that renders the QR code's finder pattern
// with rounded transitions at the corners
func RoundedFinder() func(ctx *standard.DrawContext) {
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()
		fw, fh := float64(w), float64(h)
		x, y := ctx.UpperLeft()

		lw := fw / 2
		lh := fh / 2

		ctx.SetColor(ctx.Color())

		mask := ctx.Neighbours()
		if mask&standard.NSelf != standard.NSelf {
			return
		}
		switch {
		// top right corners
		case mask == (standard.NSelf | standard.NBot | standard.NLeft):
			ctx.MoveTo(x, y)
			ctx.QuadraticTo(x+fw, y, x+fw, y+fh)
			ctx.LineTo(x, y+fh+lh)
			ctx.QuadraticTo(x, y+fh, x-lw, y+fh)
			ctx.ClosePath()
		case mask == (standard.NSelf | standard.NBot | standard.NLeft | standard.NBotLeft):
			ctx.MoveTo(x, y)
			ctx.QuadraticTo(x+fw, y, x+fw, y+fh)
			ctx.LineTo(x, y+fh)
			ctx.ClosePath()
		// top left corners
		case mask == (standard.NSelf | standard.NBot | standard.NRight):
			ctx.MoveTo(x, y+fh)
			ctx.QuadraticTo(x, y, x+fw, y)
			ctx.LineTo(x+fw+lw, y+fh)
			ctx.QuadraticTo(x+fw, y+fh, x+fw, y+fh+lh)
			ctx.ClosePath()
		case mask == (standard.NSelf | standard.NBot | standard.NRight | standard.NBotRight):
			ctx.MoveTo(x, y+fh)
			ctx.QuadraticTo(x, y, x+fw, y)
			ctx.LineTo(x+fw, y+fh)
			ctx.ClosePath()
		// bot left corners
		case mask == (standard.NSelf | standard.NTop | standard.NRight):
			ctx.MoveTo(x, y)
			ctx.QuadraticTo(x, y+fh, x+fw, y+fh)
			ctx.LineTo(x+fw+lw, y)
			ctx.QuadraticTo(x+fw, y, x+fw, y-lh)
			ctx.ClosePath()
		case mask == (standard.NSelf | standard.NTop | standard.NRight | standard.NTopRight):
			ctx.MoveTo(x, y)
			ctx.QuadraticTo(x, y+fh, x+fw, y+fh)
			ctx.LineTo(x+fw, y)
			ctx.ClosePath()
		// bot right corners
		case mask == (standard.NSelf | standard.NTop | standard.NLeft):
			ctx.MoveTo(x, y+fh)
			ctx.QuadraticTo(x+fw, y+fh, x+fw, y)
			ctx.LineTo(x, y-lh)
			ctx.QuadraticTo(x, y, x-lw, y)
			ctx.ClosePath()
		case mask == (standard.NSelf | standard.NTop | standard.NLeft | standard.NTopLeft):
			ctx.MoveTo(x, y+fh)
			ctx.QuadraticTo(x+fw, y+fh, x+fw, y)
			ctx.LineTo(x, y)
			ctx.ClosePath()
			ctx.Fill()
		default:
			ctx.DrawRectangle(x, y, fw, fh)
		}

		ctx.Fill()
	}
}

// SquareFinder just square finder
func SquareFinder() func(ctx *standard.DrawContext) {
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()

		fw0, fh0 := float64(w), float64(h)
		x0, y0 := ctx.UpperLeft()

		ctx.SetColor(ctx.Color())
		ctx.DrawRectangle(x0, y0, fw0, fh0)
		ctx.Fill()
	}
}
