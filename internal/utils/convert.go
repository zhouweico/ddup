package utils

import "strconv"

// StringToInt 将字符串转换为int，转换失败时返回0
func StringToInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

// StringToUint 将字符串转换为uint，转换失败时返回0
func StringToUint(str string) uint {
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint(val)
}
