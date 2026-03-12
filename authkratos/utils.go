package authkratos

// BooleanToNum converts boolean to int
// Returns 1 when true, 0 when false
//
// BooleanToNum 将布尔值转换成整数
// true 时返回 1，false 时返回 0
func BooleanToNum(b bool) int {
	if b {
		return 1
	}
	return 0
}
