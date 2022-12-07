package utils

import "strconv"

func Str2Uint(value string) uint {
	res, _ := strconv.ParseUint(value, 10, 64)
	return uint(res)
}
