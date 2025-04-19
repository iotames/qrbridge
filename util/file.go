package util

import (
	"fmt"

	"github.com/iotames/miniutils"
)

func getExistFile(defaultFilePath string, customFilePath string) string {
	if miniutils.IsPathExists(customFilePath) {
		return customFilePath
	}
	if miniutils.IsPathExists(defaultFilePath) {
		return defaultFilePath
	}
	return ""
}

func GetTextByFilePath(defaultFilePath string, customFilePath string) (content string, err error) {
	fpath := getExistFile(defaultFilePath, customFilePath)
	if fpath == "" {
		return "", fmt.Errorf("file not found: %s", defaultFilePath)
	}
	return miniutils.ReadFileToString(fpath)
}
