package utils

import "unicode/utf8"

// GetLongest multiples val by per
func GetLongest(seqs []string) int {
	n := 1
	for _, seq := range seqs {
		c := utf8.RuneCountInString(seq)
		if n < c {
			n = c
		}
	}
	return n
}
