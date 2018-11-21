package qrcode

import (
	"testing"

	"github.com/yeqown/go-qrcode/matrix"
)

func Test_image_draw(t *testing.T) {
	m := matrix.NewMatrix(20, 20)
	// set all 3rd column as black else be white
	for x := 0; x < m.Width(); x++ {
		m.Set(x, 3)
	}

	if err := draw("./testdata/default.jpeg", *m); err != nil {
		t.Errorf("want nil, but err: %v", err)
		t.Fail()
	}
}
