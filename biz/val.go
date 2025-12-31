package biz

import (
	"strings"
	"unicode"
)

// RemoveNonDigits 删除非数字字符
func RemoveNonDigits(input string) string {
	// 1 050
	var result strings.Builder
	for _, ch := range input {
		if unicode.IsDigit(ch) {
			result.WriteRune(ch)
		}
	}
	return result.String()
}
