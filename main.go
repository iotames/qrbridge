package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/webserver"
)

// 使用-ldflags参数在编译时设置 Version 和 DbFlag 的值
// For Linux: go build -v -ldflags "-X 'main.BuildTime=$(date +%Y-%m-%d_%H:%M)' -X 'main.Version=v1.1.0' -X 'main.DbFlag=false' "
// For Windows: go build -v -o PO转换工具.exe -trimpath -ldflags "-X 'main.BuildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%:%time:~3,2%' -X 'main.Version=v1.1.0' -X 'main.DbFlag=false' " .
var (
	BuildTime string
	Version   = "v1.2.3"
	DbFlag    = "true"
)

func main() {
	if vsion {
		fmt.Printf("QRBridge:%s, BuildTime:%s\n", Version, BuildTime)
		fmt.Println("DbFlag", DbFlag)
		return
	}
	if checkDaemon() {
		return
	}
	if DbFlag == "true" {
		CheckDb()
	}
	if Debug {
		debug()
		return
	}

	if IsPathExists("tpl/amis.html") {
		go func() {
			time.Sleep(1 * time.Second)
			err := startBrowser()
			if err != nil {
				log.Println("startBrowser Error:", err)
			}
		}()
	}

	webserver.Run(conf.WebServerPort)
}

func init() {
	err := conf.LoadEnv()
	if err != nil {
		panic(fmt.Errorf("init err(%v)", err))
	}
	parseArgs()
	initScript()
}

func CheckDb() bool {
	if conf.DbDriver == "" {
		fmt.Println("警告：数据库驱动 DB_DRIVER 未配置")
	} else {
		err := db.DbOpen(conf.DbPort, conf.DbDriver, conf.DbHost, conf.DbUsername, conf.DbPassword, conf.DbName)
		if err != nil {
			panic(err)
		}
		err = db.GetDbOpen().Ping()
		if err != nil {
			panic(err)
		}
		defer db.DbClose()
		CheckDbInit()
	}
	return true
}

func checkDaemon() bool {
	args := os.Args

	for k, v := range args {
		if v == "-d" || v == "--d" {
			Daemon = true
			args[k] = ""
		}
		if v == "stop" {
			fmt.Println("停止功能不可用")
			return true
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
		return true
	}
	return false
}
