package terminal

//
//import (
//	"fmt"
//	"image/color"
//)
//
//var (
//	White = color.RGBA{R: 255, G: 255, B: 255, A: 255}
//	Black = color.RGBA{A: 255}
//)
//
//// hexToRGBA convert hex string into color.RGBA
//func hexToRGBA(s string) color.RGBA {
//	c := color.RGBA{
//		R: 0,
//		G: 0,
//		B: 0,
//		A: 0xff,
//	}
//
//	var err error
//	switch len(s) {
//	case 7:
//		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
//	case 4:
//		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
//		// Double the hex digits:
//		c.R *= 17
//		c.G *= 17
//		c.B *= 17
//	default:
//		err = fmt.Errorf("invalid length, must be 7 or 4")
//	}
//	if err != nil {
//		panic(err)
//	}
//
//	return c
//}
