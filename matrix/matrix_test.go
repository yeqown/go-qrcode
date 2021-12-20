package matrix

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrix(t *testing.T) {
	m := New(3, 3)
	m.print()

	err := m.Set(2, 4, StateTrue)
	assert.Error(t, err)

	err = m.Set(2, 2, StateTrue)
	assert.NoError(t, err)

	v, err2 := m.Get(2, 2)
	assert.NoError(t, err2)
	assert.Equal(t, StateTrue, v)

	m.print()
	// m.Reset(2, 2)
}

func TestXOR(t *testing.T) {
	type args struct {
		s1 State
		s2 State
	}
	tests := []struct {
		name string
		args args
		want State
	}{
		{
			name: "case1",
			args: args{
				s1: StateFalse,
				s2: StateFalse,
			},
			want: StateFalse,
		},
		{
			name: "case1",
			args: args{
				s1: StateTrue,
				s2: StateFalse,
			},
			want: StateTrue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := XOR(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("XOR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_String(t *testing.T) {
	tests := []struct {
		name string
		s    State
		want string
	}{
		{
			name: "case 1",
			s:    StateFalse,
			want: "0x1",
		},
		{
			name: "case 2",
			s:    StateTrue,
			want: "0x2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("State.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Copy(t *testing.T) {
	// pre
	var (
		m1 = New(3, 3)
		m2 = New(3, 3)
	)
	m1.Set(1, 1, StateFalse)
	m1.Set(0, 0, StateFalse)
	m2.Set(1, 1, StateFalse)
	m2.Set(0, 0, StateFalse)

	// do copy
	got := m1.Copy()

	// change origin
	m1.Set(2, 2, StateTrue)

	// compare copy of the matrix
	if !reflect.DeepEqual(got, m2) {
		t.Errorf("Matrix.Copy() = %v, want %v", got, m2)
		t.Fail()
	}

	if s, err := m1.Get(2, 2); err != nil {
		t.Errorf("Matrix.Get() = %v, want %v", StateTrue, s)
		t.Fail()
	} else if s != StateTrue {
		t.Errorf("Matrix.Copy() = %v, want %v", StateTrue, s)
		t.Fail()
	}

	if s, _ := got.Get(2, 2); s != StateInit {
		t.Errorf("Matrix.Copy() = %v, want %v", StateInit, s)
		t.Fail()
	}
}

func Test_stateSliceMatched(t *testing.T) {
	type args struct {
		ss1 []State
		ss2 []State
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{
				ss1: []State{StateFalse, StateFalse, StateFalse, StateFalse, StateTrue},
				ss2: []State{StateFalse, StateFalse, StateFalse, StateFalse, StateTrue},
			},
			want: true,
		},
		{
			name: "case 0",
			args: args{
				ss1: []State{StateFalse, StateFalse, StateFalse, StateFalse, StateTrue},
				ss2: []State{StateFalse, StateTrue, StateFalse, StateFalse, StateTrue},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StateSliceMatched(tt.args.ss1, tt.args.ss2); got != tt.want {
				t.Errorf("stateSliceMatched() = %v, want %v", got, tt.want)
			}
		})
	}
}

// go test -run=NONE -bench Benchmark_Iterate -count 10 > old.txt
// after change to `m.Iterate(COLUMN, rowIteration)`
// go test -run=NONE -bench Benchmark_Iterate -count 10 > new.txt
// benchstat old.txt new.txt
func Benchmark_Iterate(b *testing.B) {
	// initialize
	size := 100
	m := New(size, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			_ = m.Set(i, j, StateTrue)
		}
	}
	b.ResetTimer()

	rowIteration := func(x, y int, s State) {
		_, _ = x, y
		_ = s
	}

	for i := 0; i < b.N; i++ {
		//m.Iterate(ROW, rowIteration)
		m.Iterate(COLUMN, rowIteration)
	}
}
