package biz

// GetDigits 获取ASCII数字串
func GetDigits(s string) string {
	// 1 050
	var result []byte
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = append(result, s[i])
		}
	}
	return string(result)
}
