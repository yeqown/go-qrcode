package terminal

//
//import (
//	"image/color"
//	"os"
//	"syscall"
//	"unsafe"
//)
//
//const defaultRatio float64 = 7.0 / 3.0 // The terminal's default cursor width/height ratio
//
//func terminalSize() (int, int, float64) {
//	var size [4]uint16
//	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL,
//		uintptr(os.Stdout.Fd()), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&size)),
//		0, 0, 0); err != 0 {
//		panic(err)
//	}
//	rows, cols, width, height := size[0], size[1], size[2], size[3]
//
//	whratio := defaultRatio
//	if width > 0 && height > 0 {
//		whratio = float64(height/rows) / float64(width/cols)
//	}
//
//	return int(cols), int(rows), whratio
//}
//
//func terminalColor(rgb color.RGBA) uint16 {
//	r := (((rgb.R * 5) + 127) / 255) * 36
//	g := (((rgb.G * 5) + 127) / 255) * 6
//	b := ((rgb.B * 5) + 127) / 255
//
//	return uint16(r+g+b) + 16 + 1 // termbox default color offset
//}
