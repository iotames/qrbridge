package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

// GetDbOpen 获取数据库连接的单例实例
var (
	once     sync.Once
	instance *sql.DB
)

func GetDbOpen() *sql.DB {
	once.Do(func() {
		if instance == nil {
			panic("请先调用DbOpen方法初始化数据库连接")
		}
	})
	return instance
}

func DbOpen(dbPort int, driverName, dbHost, dbUser, dbPassword, dbName string) error {
	var err error

	// 构建数据库连接字符串，包含schema参数
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)
	// 使用sql.Open打开数据库连接池
	instance, err = sql.Open(driverName, connStr)
	return err
}

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

// DbClose 关闭数据库连接
func DbClose() error {
	if instance != nil {
		return GetDbOpen().Close()
	}
	return nil
}
