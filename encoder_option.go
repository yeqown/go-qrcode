package qrcode

// DefaultConfig with EncMode = EncModeAuto, EcLevel = ErrorCorrectionQuart
func DefaultConfig() *Config {
	return &outputEncodingOption{
		EncMode: EncModeAuto,
		EcLevel: ErrorCorrectionQuart,
	}
}

// Config alias of outputEncodingOption.
type Config = outputEncodingOption

type outputEncodingOption struct {

	// EncMode specifies which encMode to use
	EncMode encMode

	// EcLevel specifies which ecLevel to use
	EcLevel ecLevel

	// PS: The version (which implicitly defines the byte capacity of the qrcode) is dynamically selected at runtime
}
