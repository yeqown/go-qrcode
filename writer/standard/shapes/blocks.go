package shapes

import "github.com/yeqown/go-qrcode/writer/standard"

// LiquidBlock returns a drawing function that renders QR-code-like blocks
// with smooth, organic, fluid transitions based on neighboring cell presence.
// It creates a visually connected, blob-style appearance by dynamically adjusting
// corners and sides depending on the surrounding cell mask.
func LiquidBlock() func(ctx *standard.DrawContext) {
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()
		fw, fh := float64(w), float64(h)
		x, y := ctx.UpperLeft()
		cx, cy := x+fw/2, y+fh/2
		r := fw / 2
		l := fw / 2

		ctx.SetColor(ctx.Color())

		type angleDrawer func(ctx *standard.DrawContext)
		var (
			AngTopRight angleDrawer = func(ctx *standard.DrawContext) {
				ctx.MoveTo(cx, cy+r)
				ctx.LineTo(cx-r, cy)
				ctx.LineTo(cx-r, y)
				ctx.LineTo(cx+r, y-l)
				ctx.QuadraticTo(cx+r, cy-r, x+fw+l, cy-r)
				ctx.LineTo(x+fw, cy+r)
				ctx.ClosePath()
			}
			AngTopLeft angleDrawer = func(ctx *standard.DrawContext) {
				ctx.MoveTo(cx, cy+r)
				ctx.LineTo(cx+r, cy)
				ctx.LineTo(cx+r, y-l)
				ctx.LineTo(cx-r, y-l)
				ctx.QuadraticTo(cx-r, cy-r, x-l, cy-r)
				ctx.LineTo(x-l, cy+r)
				ctx.ClosePath()
			}
			AngBotLeft angleDrawer = func(ctx *standard.DrawContext) {
				ctx.MoveTo(cx, cy-r)
				ctx.LineTo(cx+r, cy)
				ctx.LineTo(cx+r, y+fh+l)
				ctx.LineTo(cx-r, y+fh+l)
				ctx.QuadraticTo(cx-r, cy+r, x-l, cy+r)
				ctx.LineTo(x-l, cy-r)
				ctx.ClosePath()
			}
			AngBotRight angleDrawer = func(ctx *standard.DrawContext) {
				ctx.MoveTo(cx, cy-r)
				ctx.LineTo(cx-r, cy)
				ctx.LineTo(cx-r, y+fh+l)
				ctx.LineTo(cx+r, y+fh+l)
				ctx.QuadraticTo(cx+r, cy+r, x+fw+l, cy+r)
				ctx.LineTo(x+fw, cy-r)
				ctx.ClosePath()
			}
		)

		mask := ctx.Neighbours()

		drawRect := func(x, y, w, h float64) {
			ctx.DrawRectangle(x, y, w, h)
			ctx.Fill()
		}

		switch mask {
		case standard.NRight | standard.NSelf:
			drawRect(cx, cy-r, fw/2, 2*r)
		case standard.NTop | standard.NSelf:
			drawRect(cx-r, y, 2*r, fh/2)
		case standard.NLeft | standard.NSelf:
			drawRect(x, cy-r, fw/2, 2*r)
		case standard.NBot | standard.NSelf:
			drawRect(cx-r, y+fh/2, 2*r, fh/2)
		case standard.NLeft | standard.NSelf | standard.NRight:
			drawRect(x-fw/2, cy-r, 2*fw, 2*r)
		}

		if has(mask, standard.NLeft|standard.NSelf|standard.NRight) {
			drawRect(x-fw/2, cy-r, 2*fw, 2*r)
		}
		if has(mask, standard.NTop|standard.NSelf|standard.NBot) {
			drawRect(cx-r, y-fh/2, 2*r, 2*fh)
		}
		if has(mask, standard.NLeft|standard.NSelf) {
			drawRect(x, cy-r, fw/2, 2*r)
		}
		if has(mask, standard.NSelf|standard.NRight) {
			drawRect(cx, cy-r, fw/2, 2*r)
		}
		if has(mask, standard.NSelf|standard.NTop) {
			drawRect(cx-r, y, 2*r, fh/2)
		}
		if has(mask, standard.NSelf|standard.NBot) {
			drawRect(cx-r, y+fh/2, 2*r, fh/2)
		}

		if has(mask, standard.NBot|standard.NRight|standard.NSelf) && mask&standard.NBotRight == 0 {
			AngBotRight(ctx)
			ctx.Fill()
		}
		if has(mask, standard.NBot|standard.NLeft|standard.NSelf) && mask&standard.NBotLeft == 0 {
			AngBotLeft(ctx)
			ctx.Fill()
		}
		if has(mask, standard.NTop|standard.NLeft|standard.NSelf) && mask&standard.NTopLeft == 0 {
			AngTopLeft(ctx)
			ctx.Fill()
		}
		if has(mask, standard.NTop|standard.NRight|standard.NSelf) && mask&standard.NTopRight == 0 {
			AngTopRight(ctx)
			ctx.Fill()
		}

		ctx.DrawCircle(cx, cy, r)
		ctx.Fill()
	}
}

// HStripeBlock returns a drawing function that renders a QR-code-like block
// with horizontal stripe connections. The width of each stripe is determined
// by the stripeRatio parameter, which should be in the range [0.6, 1.0].
// If the given ratio is out of bounds, a default of 0.85 is used.
func HStripeBlock(stripeRatio float64) func(ctx *standard.DrawContext) {
	if stripeRatio < 0.6 || stripeRatio > 1 {
		stripeRatio = 0.85
	}
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()
		fw, fh := float64(w), float64(h)
		x, y := ctx.UpperLeft()
		cx, cy := x+fw/2, y+fh/2
		r := fw * 0.9 / 2

		ctx.SetColor(ctx.Color())

		mask := ctx.Neighbours()

		drawRect := func(x, y, w, h float64) {
			ctx.DrawRectangle(x, y, w, h)
			ctx.Fill()
		}

		ctx.DrawCircle(cx, cy, r)

		if has(mask, standard.NLeft|standard.NSelf) {
			drawRect(x, cy-r, fw/2, 2*r)
		}
		if has(mask, standard.NRight|standard.NSelf) {
			drawRect(cx, cy-r, fw/2, 2*r)
		}

		ctx.Fill()
	}
}

// VStripeBlock returns a drawing function that renders a QR-code-like block
// with vertical stripe connections. The width of each stripe is determined
// by the stripeRatio parameter, which should be in the range [0.6, 1.0].
// If the given ratio is out of bounds, a default of 0.85 is used.
func VStripeBlock(stripeRatio float64) func(ctx *standard.DrawContext) {
	if stripeRatio < 0.6 || stripeRatio > 1 {
		stripeRatio = 0.85
	}
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()
		fw, fh := float64(w), float64(h)
		x, y := ctx.UpperLeft()
		cx, cy := x+fw/2, y+fh/2
		r := fw * stripeRatio / 2

		ctx.SetColor(ctx.Color())

		mask := ctx.Neighbours()

		drawRect := func(x, y, w, h float64) {
			ctx.DrawRectangle(x, y, w, h)
			ctx.Fill()
		}
		_ = mask
		_ = drawRect

		ctx.DrawCircle(cx, cy, r)

		if has(mask, standard.NTop|standard.NSelf) {
			drawRect(cx-r, y, 2*r, fh/2)
		}
		if has(mask, standard.NBot|standard.NSelf) {
			drawRect(cx-r, cy, 2*r, fh/2)
		}

		ctx.Fill()
	}
}

// ChainBlock returns a drawing function that renders a QR-code-like block
// with a central circle and narrow "chain link" connectors extending in all four directions.
func ChainBlock() func(ctx *standard.DrawContext) {
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()
		fw, fh := float64(w), float64(h)
		x, y := ctx.UpperLeft()
		cx, cy := x+fw/2, y+fh/2
		r := fw * 0.9 / 2
		l := r * 0.2

		ctx.SetColor(ctx.Color())

		mask := ctx.Neighbours()

		drawRect := func(x, y, w, h float64) {
			ctx.DrawRectangle(x, y, w, h)
			ctx.Fill()
		}
		_ = mask
		_ = drawRect

		ctx.DrawCircle(cx, cy, r)

		if has(mask, standard.NTop|standard.NSelf) {
			drawRect(cx-l, y, 2*l, fh/2)
		}
		if has(mask, standard.NBot|standard.NSelf) {
			drawRect(cx-l, cy, 2*l, fh/2)
		}
		if has(mask, standard.NLeft|standard.NSelf) {
			drawRect(x, cy-l, fw/2, 2*l)
		}
		if has(mask, standard.NRight|standard.NSelf) {
			drawRect(cx, cy-l, fw/2, 2*l)
		}

		ctx.Fill()
	}
}

// VChainBlock returns a drawing function that renders a QR-code-like block
// with a central circle and narrow vertical "chain link" connectors to
// neighboring blocks above and below.
func VChainBlock() func(ctx *standard.DrawContext) {
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()
		fw, fh := float64(w), float64(h)
		x, y := ctx.UpperLeft()
		cx, cy := x+fw/2, y+fh/2
		r := fw * 0.85 / 2 // todo:
		l := r * 0.2

		ctx.SetColor(ctx.Color())

		mask := ctx.Neighbours()

		drawRect := func(x, y, w, h float64) {
			ctx.DrawRectangle(x, y, w, h)
			ctx.Fill()
		}
		_ = mask
		_ = drawRect

		ctx.DrawCircle(cx, cy, r)

		if has(mask, standard.NTop|standard.NSelf) {
			drawRect(cx-l, y, 2*l, fh/2)
		}
		if has(mask, standard.NBot|standard.NSelf) {
			drawRect(cx-l, cy, 2*l, fh/2)
		}

		ctx.Fill()
	}
}

// HChainBlock returns a drawing function that renders a QR-code-like block
// with a central circle and narrow horizontal "chain link" connectors to
// neighboring blocks on the left and right.
func HChainBlock() func(ctx *standard.DrawContext) {
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()
		fw, fh := float64(w), float64(h)
		x, y := ctx.UpperLeft()
		cx, cy := x+fw/2, y+fh/2
		r := fw * 0.85 / 2 // todo:
		l := r * 0.2

		ctx.SetColor(ctx.Color())

		mask := ctx.Neighbours()

		drawRect := func(x, y, w, h float64) {
			ctx.DrawRectangle(x, y, w, h)
			ctx.Fill()
		}
		_ = mask
		_ = drawRect

		ctx.DrawCircle(cx, cy, r)

		if has(mask, standard.NLeft|standard.NSelf) {
			drawRect(x, cy-l, fw/2, 2*l)
		}
		if has(mask, standard.NRight|standard.NSelf) {
			drawRect(cx, cy-l, fw/2, 2*l)
		}

		ctx.Fill()
	}
}

// SquareBlocks returns a drawing function that renders a centered square block.
// The size parameter defines the square's size relative to the available cell,
// ranging from 0.1 (10%) to 1.0 (100%). Values outside this range default to 1.0.
func SquareBlocks(size float64) func(ctx *standard.DrawContext) {
	if size < 0.1 || size > 1.0 {
		size = 1
	}
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()

		fw0, fh0 := float64(w), float64(h)
		x0, y0 := ctx.UpperLeft()
		cx, cy := x0+fw0/2, y0+fh0/2

		fw, fh := fw0*size, fh0*size
		x, y := cx-fw/2, cy-fh/2

		ctx.SetColor(ctx.Color())
		ctx.DrawRectangle(x, y, fw, fh)
		ctx.Fill()
	}
}

func CircleBlocks(size float64) func(ctx *standard.DrawContext) {
	if size < 0.1 || size > 1.0 {
		size = 1
	}
	return func(ctx *standard.DrawContext) {
		w, h := ctx.Edge()

		fw0, fh0 := float64(w), float64(h)
		x0, y0 := ctx.UpperLeft()
		cx, cy := x0+fw0/2, y0+fh0/2

		r := (fw0 / 2) * size

		ctx.SetColor(ctx.Color())
		ctx.DrawCircle(cx, cy, r)
		ctx.Fill()
	}
}
