package qrcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//func Test_rule3_refactor(t *testing.T) {
//	qrc, err := newMatrix("baidu.com google.com qq.com sina.com apple.com")
//	assert.NoError(t, err)
//	_ = qrc
//	old := rule3_backup(qrc.mat)
//	refactor := rule3(qrc.mat)
//	assert.Equal(t, old, refactor)
//}

func Benchmark_rule3(b *testing.B) {
	qrc, err := New("baidu.com google.com qq.com sina.com apple.com")
	assert.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rule3(qrc.mat)
	}
}

//func Test_rule1_refactor(t *testing.T) {
//	qrc, err := newMatrix("baidu.com google.com qq.com sina.com apple.com")
//	assert.NoError(t, err)
//	qrc.mat.Print()
//
//	old := rule1_backup(qrc.mat)
//	refactor := rule1(qrc.mat)
//	assert.Equal(t, old, refactor)
//}

func Benchmark_rule1(b *testing.B) {
	qrc, err := New("baidu.com google.com qq.com sina.com apple.com")
	assert.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rule1(qrc.mat)
	}
}

//func Test_rule4_refactor(t *testing.T) {
//	qrc, err := newMatrix("baidu.com google.com qq.com sina.com apple.com")
//	assert.NoError(t, err)
//	//qrc.mat.Print()
//
//	old := rule4_backup(qrc.mat)
//	refactor := rule4(qrc.mat)
//	assert.Equal(t, old, refactor)
//}
