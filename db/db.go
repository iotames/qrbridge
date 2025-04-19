package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

// GetDbOpen 获取数据库连接的单例实例
var (
	once     sync.Once
	instance *MyDb
)

func GetDbOpen() *MyDb {
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
	var sqldb *sql.DB
	sqldb, err = sql.Open(driverName, connStr)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}
	instance = NewMyDb(sqldb)
	return err
}

// DbClose 关闭数据库连接
func DbClose() error {
	if instance != nil {
		return GetDbOpen().Close()
	}
	return nil
}
