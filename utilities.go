package qrcode

import "github.com/yeqown/go-qrcode/v2/matrix"

// samestate judge two matrix State is same with binary semantic.
// StateFalse/StateInit only equal to StateFalse, other state are equal to each other.
func samestate(s1, s2 matrix.State) bool {
	if s1 == s2 {
		return true
	}

	switch s1 {
	case matrix.StateFalse, matrix.StateInit:
		return false
	}
	switch s2 {
	case matrix.StateFalse, matrix.StateInit:
		return false
	}

	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}

func binaryToStateSlice(s string) []matrix.State {
	var states = make([]matrix.State, 0, len(s))
	for _, c := range s {
		switch c {
		case '1':
			states = append(states, matrix.StateTrue)
		case '0':
			states = append(states, matrix.StateFalse)
		default:
			continue
		}
	}
	return states
}
