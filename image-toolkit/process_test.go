package kit_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	kit "github.com/yeqown/go-qrcode/image-toolkit"
)

func Test_Gray(t *testing.T) {
	img, err := kit.Read("testdata/test.png")
	assert.NoError(t, err)

	out := kit.Gray(img)
	assert.Equal(t, out.Bounds(), img.Bounds())
	kit.Save(out, "testdata/test_gray.png")
}

func TestBinaryzation(t *testing.T) {
	img, err := kit.Read("testdata/test.png")
	assert.NoError(t, err)

	out := kit.Binaryzation(img, 60)
	assert.Equal(t, out.Bounds(), img.Bounds())
	err = kit.Save(out, "testdata/test_binaryzation.png")
	assert.NoError(t, err)
}

func TestScale(t *testing.T) {
	img, err := kit.Read("testdata/test_binaryzation.png")
	assert.NoError(t, err)

	out := kit.Scale(img, image.Rect(0, 0, 100, 100), nil)
	assert.Equal(t, out.Bounds(), image.Rect(0, 0, 100, 100))
	err = kit.Save(out, "testdata/test_binaryzation_scale.png")
	assert.NoError(t, err)
}
