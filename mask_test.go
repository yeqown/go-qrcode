package qrcode

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yeqown/go-qrcode/v2/matrix"
)

func TestMask(t *testing.T) {
	qrc := &QRCode{
		content:     "baidu.com google.com qq.com sina.com apple.com",
		mode:        DefaultEncodingOption().EncMode,
		ecLv:        DefaultEncodingOption().EcLevel,
		needAnalyze: true,
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
	_ = debugDraw("./assets/mask_origin.jpeg", *cpyMat)

	mask0 := newMask(cpyMat, modulo0)
	_ = debugDraw("./assets/modulo0.jpeg", *mask0.mat)

	mask1 := newMask(cpyMat, modulo1)
	_ = debugDraw("./assets/modulo1.jpeg", *mask1.mat)

	mask2 := newMask(cpyMat, modulo2)
	_ = debugDraw("./assets/modulo2.jpeg", *mask2.mat)

	mask3 := newMask(cpyMat, modulo3)
	_ = debugDraw("./assets/modulo3.jpeg", *mask3.mat)

	mask4 := newMask(cpyMat, modulo4)
	_ = debugDraw("./assets/modulo4.jpeg", *mask4.mat)

	mask5 := newMask(cpyMat, modulo5)
	_ = debugDraw("./assets/modulo5.jpeg", *mask5.mat)

	mask6 := newMask(cpyMat, modulo6)
	_ = debugDraw("./assets/modulo6.jpeg", *mask6.mat)

	mask7 := newMask(cpyMat, modulo7)
	_ = debugDraw("./assets/modulo7.jpeg", *mask7.mat)
}
