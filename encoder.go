// Package qrcode ...
// encoder.go working for data encoding
package qrcode

import (
	"github.com/yeqown/go-qrcode/bitset"
)

// Encoder ... data to bit stream ...
type Encoder struct {
	dst *bitset.Bitset
	// raw input data
	data []byte
}

// Encode ...
func (e *Encoder) Encode(byts []byte) error {

}
