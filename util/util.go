package util

import "strconv"

func Str2Float64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
