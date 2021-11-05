package qrcode

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yeqown/go-qrcode/matrix"
)

func TestMask(t *testing.T) {
	qrc := &QRCode{
		content:      "baidu.com google.com qq.com sina.com apple.com",
		mode:         DefaultConfig().EncMode,
		ecLv:         DefaultConfig().EcLevel,
		needAnalyze:  true,
		outputOption: defaultOutputImageOption(),
	}
	err := qrc.init()
	require.NoError(t, err)

	var stateInitCnt int
	qrc.mat.Iterate(matrix.ROW, func(x, y int, s matrix.State) {
		if s == matrix.StateInit {
			stateInitCnt++
		}
	})
	t.Logf("all StateInit block count: %d", stateInitCnt)

	cpyMat := qrc.mat.Copy()
	_ = drawAndSaveToFile("./testdata/mask_origin.jpeg", *cpyMat, nil)

	mask0 := newMask(cpyMat, modulo0)
	_ = drawAndSaveToFile("./testdata/modulo0.jpeg", *mask0.mat, nil)

	mask1 := newMask(cpyMat, modulo1)
	_ = drawAndSaveToFile("./testdata/modulo1.jpeg", *mask1.mat, nil)

	mask2 := newMask(cpyMat, modulo2)
	_ = drawAndSaveToFile("./testdata/modulo2.jpeg", *mask2.mat, nil)

	mask3 := newMask(cpyMat, modulo3)
	_ = drawAndSaveToFile("./testdata/modulo3.jpeg", *mask3.mat, nil)

	mask4 := newMask(cpyMat, modulo4)
	_ = drawAndSaveToFile("./testdata/modulo4.jpeg", *mask4.mat, nil)

	mask5 := newMask(cpyMat, modulo5)
	_ = drawAndSaveToFile("./testdata/modulo5.jpeg", *mask5.mat, nil)

	mask6 := newMask(cpyMat, modulo6)
	_ = drawAndSaveToFile("./testdata/modulo6.jpeg", *mask6.mat, nil)

	mask7 := newMask(cpyMat, modulo7)
	_ = drawAndSaveToFile("./testdata/modulo7.jpeg", *mask7.mat, nil)
}
