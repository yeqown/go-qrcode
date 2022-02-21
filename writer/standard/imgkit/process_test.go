package imgkit_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeqown/go-qrcode/writer/standard/imgkit"
)

func Test_Gray(t *testing.T) {
	t.Skipf("need human to check")

	img, err := imgkit.Read("testdata/test.png")
	assert.NoError(t, err)

	out := imgkit.Gray(img)
	assert.Equal(t, out.Bounds(), img.Bounds())
	imgkit.Save(out, "testdata/test_gray.png")
}

func TestBinaryzation(t *testing.T) {
	t.Skipf("need human to check")

	img, err := imgkit.Read("testdata/test.png")
	assert.NoError(t, err)

	out := imgkit.Binaryzation(img, 60)
	assert.Equal(t, out.Bounds(), img.Bounds())
	err = imgkit.Save(out, "testdata/test_binaryzation.png")
	assert.NoError(t, err)
}

func TestScale(t *testing.T) {
	t.Skipf("need human to check")

	img, err := imgkit.Read("testdata/test_binaryzation.png")
	assert.NoError(t, err)

	out := imgkit.Scale(img, image.Rect(0, 0, 100, 100), nil)
	assert.Equal(t, out.Bounds(), image.Rect(0, 0, 100, 100))
	err = imgkit.Save(out, "testdata/test_binaryzation_scale.png")
	assert.NoError(t, err)
}
