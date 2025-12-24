package main

import (
	"fmt"
	"log"

	"time"

	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/webserver"
)

// 使用-ldflags参数在编译时设置 Version 和 DbFlag 的值
// For Linux: go build -v -ldflags "-X 'main.BuildTime=$(date +%Y-%m-%d_%H:%M)' -X 'main.Version=v1.1.0' -X 'main.DbFlag=false' "
// For Windows: go build -v -o PO转换工具.exe -trimpath -ldflags "-X 'main.BuildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%:%time:~3,2%' -X 'main.Version=v1.1.0' -X 'main.DbFlag=false' " .
var (
	BuildTime string
	Version   = "v1.3.0"
	DbFlag    = "true"
)

func main() {
	if Debug {
		debug()
		return
	}
	if vsion {
		fmt.Printf("QRBridge:%s, BuildTime:%s\n", Version, BuildTime)
		fmt.Println("DbFlag", DbFlag)
		return
	}
	if Inputtpl != "" && Inputfile != "" {
		err := inputfileTransform()
		if err != nil {
			panic(err)
		}
		return
	}
	if checkDaemon() {
		return
	}
	if DbFlag == "true" {
		CheckDb()
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
