package common

import (
	"strconv"
	"strings"
)

func Atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func Atoi8(s string) int8 {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 8)
	return int8(i)
}

func Atoi16(s string) int16 {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 16)
	return int16(i)
}

func Atoi32(s string) int32 {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return int32(i)
}

func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return i
}
