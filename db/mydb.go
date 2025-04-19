package db

import (
	"database/sql"
	"log"
	"time"
)

// MyDb 是一个包装了sql.DB的结构体，用于记录SQL查询
type MyDb struct {
	*sql.DB
}

// 创建一个新的MyDb实例
func NewMyDb(db *sql.DB) *MyDb {
	return &MyDb{db}
}

// Exec 重写Exec方法以记录SQL查询
func (db *MyDb) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	// log.Printf("执行SQL: %s 参数: %v", query, args)
	result, err := db.DB.Exec(query, args...)
	log.Printf("SQL执行完成，耗时: %v", time.Since(start))
	return result, err
}

// Query 重写Query方法以记录SQL查询
func (db *MyDb) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	// log.Printf("查询SQL: %s 参数: %v", query, args)
	rows, err := db.DB.Query(query, args...)
	log.Printf("SQL查询完成，耗时: %v", time.Since(start))
	return rows, err
}

// QueryRow 重写QueryRow方法以记录SQL查询
func (db *MyDb) QueryRow(query string, args ...interface{}) *sql.Row {
	// log.Printf("查询单行SQL: %s 参数: %v", query, args)
	return db.DB.QueryRow(query, args...)
}
