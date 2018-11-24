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

	if err := m.Reset(2, 4); err == nil {
		t.Errorf("want err, got nil")
		t.Fail()
	}

	if err := m.Set(2, 2, StateInit); err != nil {
		t.Errorf("want nil, got err: %v", err)
		t.Fail()
	}

	if v, err := m.Get(2, 2); err != nil || v != StateInit {
		t.Errorf("want true, got %v and err: %v", v, err)
		t.Fail()
	}

	m.print()
	m.Reset(2, 2)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
