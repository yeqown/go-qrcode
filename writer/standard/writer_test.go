package standard

import (
	"crypto/md5"
	"encoding/hex"
	"image/png"
	"io"
	"os"
	"testing"

	"github.com/yeqown/go-qrcode/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	qrc, err := qrcode.New("cost 3+ days to read docs and coding")
	require.NoError(t, err)

	w, err := New("./testdata/qrtest.jpeg")
	require.NoError(t, err)

	// save file
	err = qrc.Save(w)
	require.NoError(t, err)

	// check file existed
	_, err = os.Stat("./testdata/qrtest.jpeg")
	if err != nil {
		t.Errorf("could not find image file: %v", err)
		t.Fail()
	}
}

func Test_New_WithOutputOption_bg_fg_width(t *testing.T) {
	qrc, err := qrcode.New("Test_New_WithOutputOption_bg_fg_width")
	require.NoError(t, err)

	w, err := New("./testdata/qrtest_fg_bg.jpeg",
		WithBgColorRGBHex("#b8de6f"),
		WithFgColorRGBHex("#01c5c4"),
		WithQRWidth(20),
	)
	require.NoError(t, err)

	// save file
	err = qrc.Save(w)
	require.NoError(t, err)
}

func Test_New_WithOutputOption_Logo(t *testing.T) {
	qrcode.SetDebugMode()

	qrc, err := qrcode.New("Test_New_WithOutputOption_Logo")
	require.NoError(t, err)

	w, err := New("./testdata/qrtest_logo.jpeg",
		WithBgColorRGBHex("#b8de6f"),
		WithFgColorRGBHex("#f1e189"),
		WithLogoImageFileJPEG("./testdata/logo.jpeg"),
		//WithLogoImageFilePNG("./testdata/logo.png"), // png required
	)
	require.NoError(t, err)

	// save file
	err = qrc.Save(w)
	require.NoError(t, err)
}

func Test_New_WithOutputOption_Shape(t *testing.T) {
	qrc, err := qrcode.New("Test_New_WithOutputOption_Shape")
	require.NoError(t, err)

	w, err := New("./testdata/qrtest_circle.jpeg",
		WithBgColorRGBHex("#b8de6f"),
		WithFgColorRGBHex("#f1e189"),
		WithCircleShape(),
	)
	require.NoError(t, err)

	err = qrc.Save(w)
	require.NoError(t, err)
}

func Test_New_WithBorderWidth(t *testing.T) {
	qrc, err := qrcode.New("Test_New_WithOutputOption_Shape")
	require.NoError(t, err)

	w, err := New("./testdata/qrtest_border_width.jpeg", WithBorderWidth(10, 20, 30, 40))

	// save file
	err = qrc.Save(w)
	require.NoError(t, err)
}

// Test_Issue40
// https://github.com/yeqown/go-qrcode/issues/40
func Test_Issue40(t *testing.T) {
	qrc, err := qrcode.New("https://yeqown.xyzom/")
	require.NoError(t, err)
	w1, err := New("./testdata/issue40_1.png")
	require.NoError(t, err)
	err = qrc.Save(w1)
	require.NoError(t, err)

	w2, err := New("./testdata/issue40_2.png")
	require.NoError(t, err)
	err = qrc.Save(w2)
	require.NoError(t, err)

	h1, err := hashFile("./testdata/issue40_1.png")
	require.NoError(t, err)
	h2, err := hashFile("./testdata/issue40_2.png")
	require.NoError(t, err)
	t.Logf("hash1=%s, hash2=%s", h1, h2)
	assert.Equal(t, h1, h2)
}

func hashFile(filename string) (string, error) {
	h := md5.New()

	fd1, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(fd1)
	if err != nil {
		return "", err
	}
	if _, err = h.Write(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func statImage(filename string) (w, h int, err error) {
	fd, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}

	img, err := png.Decode(fd)
	if err != nil {
		return 0, 0, err
	}

	rect := img.Bounds()
	w, h = rect.Dx(), rect.Dy()
	return
}

func Test_Attribute(t *testing.T) {
	qrc, err := qrcode.New("https://yeqown.xyzom")
	require.NoError(t, err)

	w, err := New("./testdata/attr.png",
		WithBuiltinImageEncoder(PNG_FORMAT),
		WithQRWidth(13),
		WithBorderWidth(1, 2, 3, 4),
	)
	require.NoError(t, err)

	attr := w.Attribute(qrc.Dimension())
	t.Logf("attr: %+v", attr)

	err = qrc.Save(w)
	require.NoError(t, err)

	width, height, err := statImage("./testdata/attr.png")
	require.NoError(t, err)
	assert.Equal(t, width, attr.W)
	assert.Equal(t, height, attr.H)
}

//
//func Test_image_draw(t *testing.T) {
//	m := new(qrcode.Matrix)
//	// set all 3rd column as black else be white
//	for x := 0; x < m.Width(); x++ {
//		_ = m.Set(x, 3, matrix.StateTrue)
//	}
//
//	fd, err := os.Create("./testdata/default.jpeg")
//	require.NoError(t, err)
//	err = drawTo(fd, *m, nil)
//	require.NoError(t, err)
//}

func Test_writer_WithBgTransparent(t *testing.T) {
	qrc, err := qrcode.New("https://yeqown.xyzom")
	require.NoError(t, err)

	w, err := New("./testdata/transparent.png",
		WithBuiltinImageEncoder(PNG_FORMAT),
		WithBorderWidth(20),
		WithBgTransparent(),
	)
	require.NoError(t, err)

	err = qrc.Save(w)
	assert.NoError(t, err)
}
