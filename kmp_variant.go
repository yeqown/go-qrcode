package qrcode

import "github.com/yeqown/go-qrcode/v2/matrix"

// kmp is variant of kmp algorithm to count the pattern been in
// src slice.
// TODO(@yeqown): implement this in generic way.
func kmp(src, pattern []matrix.State, next []int) (count int) {
	if next == nil {
		next = kmpGetNext(pattern)
	}
	slen := len(src)
	plen := len(pattern)
	i := 0 // cursor of src
	j := 0 // cursor of pattern

loop:
	for i < slen && j < plen {
		if j == -1 || src[i] == pattern[j] {
			i++
			j++
		} else {
			j = next[j]
		}
	}

	if j == plen {
		if i-j >= 0 {
			count++
		}

		// reset cursor to count duplicate pattern.
		// such as: "aaaa" and "aa", we want 3 rather than 2.
		i -= plen - 1
		j = 0

		goto loop
	}

	return count
}

func kmpGetNext(pattern []matrix.State) []int {
	fail := make([]int, len(pattern))
	fail[0] = -1

	j := 0
	k := -1

	for j < len(pattern)-1 {
		if k == -1 || pattern[j] == pattern[k] {
			k++
			j++
			fail[j] = k
		} else {
			k = fail[k]
		}
	}

	return fail
}
