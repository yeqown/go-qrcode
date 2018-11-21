package matrix

import (
	"testing"
)

func TestMatrix(t *testing.T) {
	m := NewMatrix(3, 3)
	m.print()

	if err := m.Set(2, 4); err == nil {
		t.Errorf("want err, got nil")
		t.Fail()
	}

	if err := m.Reset(2, 4); err == nil {
		t.Errorf("want err, got nil")
		t.Fail()
	}

	if err := m.Set(2, 2); err != nil {
		t.Errorf("want nil, got err: %v", err)
		t.Fail()
	}

	if v, err := m.Get(2, 2); err != nil || v != true {
		t.Errorf("want true, got %v and err: %v", v, err)
		t.Fail()
	}

	m.print()
	m.Reset(2, 2)
}
