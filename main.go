package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/webserver"
)

func main() {
	args := os.Args

	for k, v := range args {
		if v == "-d" || v == "--d" {
			Daemon = true
			args[k] = ""
		}
		if v == "stop" {
			fmt.Println("删除功能不可用")
			return
		}

	}
	if Daemon {
		fmt.Println("启动守护进程，后台运行程序")
		var newArgs []string
		if len(args) > 1 {
			newArgs = args[1:]
		}
		cmd := exec.Command(args[0], newArgs...)
		cmd.Env = os.Environ()
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		fmt.Println("守护进程启动成功，进程号：", cmd.Process.Pid)
		return
	}

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
