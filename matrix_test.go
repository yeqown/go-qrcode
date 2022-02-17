package qrcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrix(t *testing.T) {
	m := newMatrix(3, 3)
	m.print()

	err := m.set(2, 4, QRValue_DATA_V1)
	assert.Error(t, err)

	err = m.set(2, 2, QRValue_DATA_V1)
	assert.NoError(t, err)

	v, err2 := m.at(2, 2)
	assert.NoError(t, err2)
	assert.Equal(t, QRValue_DATA_V1, v)

	m.print()
	// m.Reset(2, 2)
}

func TestMatrix_Copy(t *testing.T) {
	// pre
	var (
		m1 = newMatrix(3, 3)
		m2 = newMatrix(3, 3)
	)
	_ = m1.set(1, 1, QRValue_DATA_V0)
	_ = m1.set(0, 0, QRValue_DATA_V0)
	_ = m2.set(1, 1, QRValue_DATA_V0)
	_ = m2.set(0, 0, QRValue_DATA_V0)

	// do copy
	got := m1.Copy()

	// change origin
	_ = m1.set(2, 2, QRValue_DATA_V1)
	assert.Equal(t, m2, got)

	s, err := m1.at(2, 2)
	assert.NoError(t, err)
	assert.Equal(t, QRValue_DATA_V1, s)

	s, err = got.at(2, 2)
	assert.NoError(t, err)
	assert.NotEqual(t, QRValue_DATA_V1, s)
	assert.Equal(t, QRValue_INIT_V0, s)
}

//func Test_stateSliceMatched(t *testing.T) {
//	type args struct {
//		ss1 []qrtype
//		ss2 []qrtype
//	}
//	tests := []struct {
//		name string
//		args args
//		want qrbool
//	}{
//		{
//			name: "case 0",
//			args: args{
//				ss1: []qrtype{QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V1},
//				ss2: []qrtype{QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V1},
//			},
//			want: true,
//		},
//		{
//			name: "case 0",
//			args: args{
//				ss1: []qrtype{QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V1},
//				ss2: []qrtype{QRValue_DATA_V0, QRValue_DATA_V1, QRValue_DATA_V0, QRValue_DATA_V0, QRValue_DATA_V1},
//			},
//			want: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := StateSliceMatched(tt.args.ss1, tt.args.ss2); got != tt.want {
//				t.Errorf("stateSliceMatched() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

// go test -run=NONE -bench Benchmark_Iterate -count 10 > old.txt
// after change to `m.iter(IterDirection_COLUMN, rowIteration)`
// go test -run=NONE -bench Benchmark_Iterate -count 10 > new.txt
// benchstat old.txt new.txt
func Benchmark_Iterate(b *testing.B) {
	// initialize
	size := 100
	m := newMatrix(size, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			_ = m.set(i, j, QRValue_DATA_V1)
		}
	}
	b.ResetTimer()

	rowIteration := func(x, y int, s qrvalue) {
		_, _ = x, y
		_ = s
	}

	for i := 0; i < b.N; i++ {
		//m.iter(IterDirection_ROW, rowIteration)
		m.iter(IterDirection_COLUMN, rowIteration)
	}
}

func Test_Matrix_RowAndCol(t *testing.T) {
	m := newMatrix(3, 3)
	for i := 0; i < 3; i++ {
		_ = m.set(i, 0, QRValue_DATA_V1)
		_ = m.set(i, i, QRValue_DATA_V1)
	}

	// 1 1 1
	// 0 1 0
	// 0 0 1

	tests := []struct {
		name  string
		rowed bool
		cur   int
		want  []qrvalue
	}{
		{
			name:  "row[1]",
			rowed: true,
			cur:   1,
			want:  []qrvalue{QRValue_INIT_V0, QRValue_DATA_V1, QRValue_INIT_V0},
		},
		{
			name:  "col[2]",
			rowed: false,
			cur:   2,
			want:  []qrvalue{QRValue_DATA_V1, QRValue_INIT_V0, QRValue_DATA_V1},
		},
		{
			name:  "row out",
			rowed: true,
			cur:   4,
			want:  nil,
		},
		{
			name:  "col out",
			rowed: false,
			cur:   4,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []qrvalue
			switch tt.rowed {
			case true:
				got = m.Row(tt.cur)
			default:
				got = m.Col(tt.cur)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
