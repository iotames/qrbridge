package main

import (
	"github.com/iotames/qrbridge/db"
)

func dbInit() {
	defer db.DbClose()
	err := db.ExecSqlBySqlFile("sql/init.sql")
	if err != nil {
		panic(err)
	}
}
