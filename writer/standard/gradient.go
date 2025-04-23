package standard

import (
	"image"
	"image/color"
	draw2 "image/draw"
	"math"
	"sort"
)

// ColorStop represents a single color stop in a gradient.
// T defines the position along the gradient line, ranging from 0.0 to 1.0.
// Color defines the color at this stop.
type ColorStop struct {
	T     float64 // from 0.0 to 1.0.
	Color color.RGBA
}

// LinearGradient defines a linear gradient with angle and color stops.
// The gradient progresses in the direction of the given angle (in degrees).
// Angle is interpreted as: 0 - right, 90 - up, 180 - left, 270 - down.
type LinearGradient struct {
	Stops []ColorStop // Ordered list of color stops along the gradient
	Angle float64     // Gradient angle in degrees
}

// NewGradient creates a new LinearGradient with the specified angle (in degrees) and color stops.
// The stops are sorted in ascending order of T.
func NewGradient(angle float64, stops ...ColorStop) *LinearGradient {
	sort.Slice(stops, func(i, j int) bool {
		return stops[i].T < stops[j].T
	})
	return &LinearGradient{Stops: stops, Angle: angle}

}

// applyGradient applies the linear gradient to all pixels in the image that match the given foreground color.
func (g *LinearGradient) applyGradient(img image.Image, fgColor color.RGBA) *image.RGBA {
	// Convert angle to radians and compute gradient direction vector
	angleRad := g.Angle * math.Pi / 180.0
	dx := math.Cos(angleRad)
	dy := -math.Sin(angleRad)

	bounds := img.Bounds()
	xmin, xmax := float64(bounds.Min.X), float64(bounds.Max.X)
	ymin, ymax := float64(bounds.Min.Y), float64(bounds.Max.Y)

	// Get all 4 corners of the image
	corners := [4][2]float64{
		{xmin, ymin},
		{xmin, ymax},
		{xmax, ymin},
		{xmax, ymax},
	}

	// Compute min and max projection of corners on the gradient axis
	minProj, maxProj := math.Inf(1), math.Inf(-1)
	for _, p := range corners {
		proj := p[0]*dx + p[1]*dy
		if proj < minProj {
			minProj = proj
		}
		if proj > maxProj {
			maxProj = proj
		}
	}

	// Prepare output image
	out := image.NewRGBA(bounds)
	draw2.Draw(out, bounds, img, bounds.Min, draw2.Src)

	// Replace foreGround pixels with interpolated gradient colors
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			if c == fgColor {
				// Project pixel onto gradient axis
				proj := float64(x)*dx + float64(y)*dy
				// Normalize to [0, 1]
				t := (proj - minProj) / (maxProj - minProj)
				// Set pixel color from gradient
				out.Set(x, y, interpolateColor(g.Stops, t))
			}
		}
	}
	return out
}

// interpolateColor returns a color interpolated from the gradient stops based on position t.
func interpolateColor(stops []ColorStop, t float64) color.RGBA {
	if t <= stops[0].T {
		return stops[0].Color
	}
	if t >= stops[len(stops)-1].T {
		return stops[len(stops)-1].Color
	}
	for i := 0; i < len(stops)-1; i++ {
		start := stops[i]
		end := stops[i+1]
		if t >= start.T && t <= end.T {
			// Linear interpolation between two stops
			factor := (t - start.T) / (end.T - start.T)
			return blendColors(start.Color, end.Color, factor)
		}
	}
	return color.RGBA{A: 255} // fallback (should not happen)
}

// blendColors returns the color interpolated between c1 and c2 using t in [0.0 - 1.0].
func blendColors(c1, c2 color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c1.R)*(1-t) + float64(c2.R)*t),
		G: uint8(float64(c1.G)*(1-t) + float64(c2.G)*t),
		B: uint8(float64(c1.B)*(1-t) + float64(c2.B)*t),
		A: 255,
	}
}
