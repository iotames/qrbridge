package db

import (
	"fmt"
	"os"
)

// IsTableExist 检查表是否存在
func IsTableExist(tableName string) (bool, error) {
	// 查询系统表以检查表是否存在
	query := `SELECT EXISTS (
		SELECT 1 
		FROM information_schema.tables 
		WHERE table_name = $1
	)`
	var exists bool
	err := GetDbOpen().QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("查询表是否存在失败: %v", err)
	}
	return exists, err
}

func ExecSqlBySqlFile(sqlFile string) error {
	// 读取SQL文件内容
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v", err)
	}

	// 执行SQL文本
	return ExecSqlText(string(sqlBytes))
}
