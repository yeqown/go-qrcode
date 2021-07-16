package qrcode

func defaultOutputEncoderOption() *outputEncodingOptions {
	return &outputEncodingOptions{
		encMode: EncModeAuto,
		ecLevel: ErrorCorrectionQuart,
	}
}

type outputEncodingOptions struct {

	// encMode specifies which encMode to use
	encMode encMode

	// ecLevel specifies which ecLevel to use
	ecLevel ecLevel

	// PS: The version (which implicitly defines the byte capacity of the qrcode) is dynamically selected at runtime
}
