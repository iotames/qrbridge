package main

import (
	"embed"
	"fmt"

	"github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/dbtable"
)

//go:embed sql/*
var sqlFS embed.FS

// dbInit 初始化数据库
func dbInit() {
	initSQL, err := sqlFS.ReadFile("sql/init.sql")
	if err != nil {
		panic(err)
	}
	err = db.ExecSqlText(string(initSQL))
	if err != nil {
		panic(err)
	}
}

// CheckDbInit 检查数据库是否初始化
// 检查qrcode表是否存在, 如果不存在，则调用dbInit()初始化数据库
func CheckDbInit() {
	qrcode := dbtable.Qrcode{}
	exist, err := db.IsTableExist(qrcode.TableName())
	if err != nil {
		panic(err)
	}
	if !exist {
		fmt.Println("数据库未初始化，正在初始化数据库...")
		dbInit()
	}
}
