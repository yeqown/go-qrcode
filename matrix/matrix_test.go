package matrix

import (
	"image/color"
	"reflect"
	"testing"
)

func TestMatrix(t *testing.T) {
	m := New(3, 3)
	m.print()

	if err := m.Set(2, 4, ZERO); err == nil {
		t.Errorf("want err, got nil")
		t.Fail()
	}

	// if err := m.Reset(2, 4); err == nil {
	// 	t.Errorf("want err, got nil")
	// 	t.Fail()
	// }

	if err := m.Set(2, 2, StateInit); err != nil {
		t.Errorf("want nil, got err: %v", err)
		t.Fail()
	}

	if v, err := m.Get(2, 2); err != nil || v != StateInit {
		t.Errorf("want true, got %v and err: %v", v, err)
		t.Fail()
	}

	m.print()
	// m.Reset(2, 2)
}

func Test_loadGray16(t *testing.T) {
	type args struct {
		v State
	}
	tests := []struct {
		name string
		args args
		want color.Gray16
	}{
		{
			name: "case 1",
			args: args{v: StateFalse},
			want: color.Gray16{Y: uint16(StateFalse)},
		},
		{
			name: "case 2",
			args: args{v: StateInit},
			want: color.Gray16{Y: uint16(StateInit)},
		},
		{
			name: "case 3",
			args: args{v: StateTrue},
			want: color.Gray16{Y: uint16(StateTrue)},
		},
		{
			name: "case 4",
			args: args{v: StateFormat},
			want: color.Gray16{Y: uint16(StateFormat)},
		},
		{
			name: "case 5",
			args: args{v: StateVersion},
			want: color.Gray16{Y: uint16(StateVersion)},
		},
		{
			name: "case 6",
			args: args{v: State(0x6767)},
			want: color.Gray16{Y: uint16(0x6767)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadGray16(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadGray16() = %v, want %v", got, tt.want)
			}
		})
	}
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
			want: "0xFFFF",
		},
		{
			name: "case 2",
			s:    StateTrue,
			want: "0x0",
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

	// cmpare copy of the matrix
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
}
