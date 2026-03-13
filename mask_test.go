package qrcode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMask(t *testing.T) {
	qrc := &QRCode{
		sourceText:     "baidu.com google.com qq.com sina.com apple.com",
		encodingOption: DefaultEncodingOption(),
	}
	err := qrc.init()
	require.NoError(t, err)

	var stateInitCnt int
	qrc.mat.iter(IterDirection_COLUMN, func(x, y int, s qrvalue) {
		if s.qrtype() == QRType_INIT {
			stateInitCnt++
		}
	})
	t.Logf("all QRType_INIT block count: %d", stateInitCnt)

	SetDebugMode()
	cpyMat := qrc.mat.Copy()
	_ = debugDraw("./testdata/mask_origin.jpeg", *cpyMat)

	mask0 := newMask(cpyMat, modulo0)
	_ = debugDraw("./testdata/modulo0.jpeg", *mask0.mat)

	mask1 := newMask(cpyMat, modulo1)
	_ = debugDraw("./testdata/modulo1.jpeg", *mask1.mat)

	mask2 := newMask(cpyMat, modulo2)
	_ = debugDraw("./testdata/modulo2.jpeg", *mask2.mat)

	mask3 := newMask(cpyMat, modulo3)
	_ = debugDraw("./testdata/modulo3.jpeg", *mask3.mat)

	mask4 := newMask(cpyMat, modulo4)
	_ = debugDraw("./testdata/modulo4.jpeg", *mask4.mat)

	mask5 := newMask(cpyMat, modulo5)
	_ = debugDraw("./testdata/modulo5.jpeg", *mask5.mat)

	mask6 := newMask(cpyMat, modulo6)
	_ = debugDraw("./testdata/modulo6.jpeg", *mask6.mat)

	mask7 := newMask(cpyMat, modulo7)
	_ = debugDraw("./testdata/modulo7.jpeg", *mask7.mat)
}
