package utils

import (
	"strconv"
	"strings"
)

// Percentalize multiples val by per
func Percentalize(val int, per string) int {
	sh := strings.Replace(per, "%", "", -1)
	i, _ := strconv.Atoi(sh)
	return int(float64(val) * (float64(i) * 0.01))
}
