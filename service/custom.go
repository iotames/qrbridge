package service

import (
	"path/filepath"

	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/util"
)

func GetSqlText(fpath string) (string, error) {
	defaultFilePath := filepath.Join("sql", fpath)
	customDirPath := filepath.Join(conf.CustomDir, "sql", fpath)
	return util.GetTextByFilePath(defaultFilePath, customDirPath)
}
