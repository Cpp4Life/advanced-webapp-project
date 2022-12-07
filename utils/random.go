package utils

import "github.com/thanhpk/randstr"

func GenerateRandomNumber(size int) string {
	return randstr.String(size, "0123456789")
}
