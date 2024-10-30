package utils

import "strconv"

// StringToInt 将字符串转换为int，转换失败时返回0
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// StringToUint 将字符串转换为uint，转换失败时返回0
func StringToUint(str string) uint {
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint(val)
}
