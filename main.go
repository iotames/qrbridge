package main

import (
	"fmt"

	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/webserver"
)

func main() {
	defer db.DbClose()
	if Dbinit {
		dbInit()
		return
	} else {
		CheckDbInit()
	}
	webserver.Run(fmt.Sprintf(":%d", conf.WebServerPort))
}

func init() {
	err := conf.LoadEnv()
	if err != nil {
		panic(fmt.Errorf("init err(%v)", err))
	}
	err = db.DbOpen(conf.DbPort, conf.DbDriver, conf.DbHost, conf.DbUsername, conf.DbPassword, conf.DbName)
	if err != nil {
		panic(err)
	}
	err = db.GetDbOpen().Ping()
	if err != nil {
		panic(err)
	}
	parseArgs()
}
