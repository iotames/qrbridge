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

func ExecSqlText(sqlText string) error {
	// 执行SQL文本
	_, err := GetDbOpen().Exec(sqlText)
	return err

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

// ExecSqlWithTransaction 在事务中执行多条SQL语句
func ExecSqlWithTransaction(sqlStatements []string) error {
	// 开始事务
	tx, err := GetDbOpen().Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}

	// 确保函数结束时要么提交要么回滚事务
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	// 执行每条SQL语句
	for _, sql := range sqlStatements {
		_, err := tx.Exec(sql)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("执行SQL失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// DbClose 关闭数据库连接
func DbClose() error {
	if instance != nil {
		return GetDbOpen().Close()
	}
	return nil
}
