package qrcode

type EncodeOption interface {
	apply(option *encodingOption)
}

// DefaultEncodingOption with EncMode = EncModeAuto, EcLevel = ErrorCorrectionQuart
func DefaultEncodingOption() *encodingOption {
	return &encodingOption{
		EncMode: EncModeAuto,
		EcLevel: ErrorCorrectionQuart,
	}
}

type encodingOption struct {
	// EncMode specifies which encMode to use
	EncMode encMode

	// EcLevel specifies which ecLevel to use
	EcLevel ecLevel

	// PS: The version (which implicitly defines the byte capacity of the qrcode) is dynamically selected at runtime
}

type fnEncodingOption struct {
	fn func(*encodingOption)
}

func (f fnEncodingOption) apply(option *encodingOption) {
	f.fn(option)
}

func newFnEncodingOption(fn func(*encodingOption)) fnEncodingOption {
	return fnEncodingOption{fn: fn}
}

// WithEncodingMode sets the encoding mode.
func WithEncodingMode(mode encMode) EncodeOption {
	return newFnEncodingOption(func(option *encodingOption) {
		if name := getEncModeName(mode); name == "" {
			return
		}

		option.EncMode = mode
	})
}

// WithErrorCorrectionLevel sets the error correction level.
func WithErrorCorrectionLevel(ecLevel ecLevel) EncodeOption {
	return newFnEncodingOption(func(option *encodingOption) {
		if ecLevel < ErrorCorrectionLow || ecLevel > ErrorCorrectionHighest {
			return
		}

		option.EcLevel = ecLevel
	})
}
