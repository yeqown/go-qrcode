package version

import (
	"github.com/yeqown/go-qrcode/mask"
	"github.com/yeqown/go-qrcode/matrix"
)

// RecoveryLevel ...
type RecoveryLevel int

const (
	// Low :Level L: 7% error recovery.
	Low RecoveryLevel = iota
	// Medium :Level M: 15% error recovery. Good default choice.
	Medium
	// High :Level Q: 25% error recovery.
	High
	// Highest :Level H: 30% error recovery.
	Highest
)

// QRVersion Version ...
type QRVersion interface {
	Size() int
	// init method to call by caller
	Init(RecoveryLevel) error
	// mask of the data of version
	Mask() mask.Mask
	// mark the data of version
	Data() *matrix.Matrix
	// mark the recovery level of the version
	Level() RecoveryLevel
}

// Analyze the text, and decide which version should be choose
func Analyze(text string) QRVersion {
	return &Version1{}
}
