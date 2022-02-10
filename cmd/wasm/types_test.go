package main

import (
	"encoding/json"
	"testing"
)

func Test_genOption(t *testing.T) {
	o := genOption{
		encodeOption: encodeOption{
			version: 0,
			mode:    0,
			ecLevel: "",
		},
		outputOption: outputOption{
			bgColor:       "",
			bgTransparent: false,
			qrColor:       "",
			qrWidth:       0,
			circleShape:   false,
			imageEncoder:  "",
			margin:        0,
		},
	}

	byts, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(byts))
}
