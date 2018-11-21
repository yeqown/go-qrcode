package version

// RecoveryLevel ...
type RecoveryLevel int

const (
	// L :Level L: 7% error recovery.
	L RecoveryLevel = iota
	// M :Level M: 15% error recovery. Good default choice.
	M
	// Q :Level Q: 25% error recovery.
	Q
	// H :Level H: 30% error recovery.
	H
)

// QRVersion Version ...
type QRVersion interface {
	// Version func return int (1-40)
	Version() int
	// Dimension func
	Dimension() int
	// init method to call by caller
	Init(RecoveryLevel) error
	// mask of the data of version
	// Mask() mask.Mask

	// mark the data of version
	// Data() *matrix.Matrix

	// mark the recovery level of the version
	Level() RecoveryLevel

	// Cap decide by Version and RecoveryLevel
	Cap() [4]int
}

// Analyze the text, and decide which version should be choose
// ref to: http://muyuchengfeng.xyz/%E4%BA%8C%E7%BB%B4%E7%A0%81-%E5%AD%97%E7%AC%A6%E5%AE%B9%E9%87%8F%E8%A1%A8/
func Analyze(text string) QRVersion {
	return &Version1{}
}
