package sql

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/util"
)

//go:embed *.sql
var sqlFS embed.FS

// getSqlText 获取sql文本
// TODO 优先从自定义目录中读取sql文本，如果不存在，则从默认目录中读取sql文本，如果都不存在，则从内嵌的文件中读取sql文本
func getSqlText(fpath string) (string, error) {
	defaultFilePath := filepath.Join("sql", fpath)
	customDirPath := filepath.Join(conf.CustomDir, "sql", fpath)
	return util.GetTextByFilePath(defaultFilePath, customDirPath)
}

// GetSQL 获取sql文本
// TODO replaceAny 是为了替换sql文本中的占位符, 占位符的格式为 ---%s---.只支持字符串类型的占位符
func GetSQL(fpath string, replaceAny ...any) (string, error) {
	sqlTxt, err := getSqlText(fpath)
	if err != nil {
		return "", err
	}
	sqlTxt = strings.ReplaceAll(sqlTxt, "---%s---", "%s")
	return fmt.Sprintf(sqlTxt, replaceAny...), nil
}

func LsDir() []string {
	entries, err := sqlFS.ReadDir(".")
	if err != nil {
		panic(err)
	}
	var filenames []string
	for _, entry := range entries {
		filenames = append(filenames, entry.Name())
		if entry.IsDir() {
			// 读取子目录中的文件
			subEntries, err := sqlFS.ReadDir(entry.Name())
			if err != nil {
				panic(err)
			}
			// 将子目录中的文件名添加到列表中
			for _, subEntry := range subEntries {
				filenames = append(filenames, entry.Name()+"/"+subEntry.Name())
			}
		}
	}
	return filenames
}
