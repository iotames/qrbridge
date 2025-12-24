package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/iotames/qrbridge/biz"
	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/dbtable"
	"github.com/iotames/qrbridge/hotswap"
	"github.com/iotames/qrbridge/sql"
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

func initScript() {
	sqldir := hotswap.NewScriptDir(sql.GetSqlFs(), conf.CustomDir, "sql")
	hotswap.GetScriptDir(sqldir)
}

func startBrowser() error {
	return StartBrowserByUrl(fmt.Sprintf("http://127.0.0.1:%d", conf.WebServerPort))
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

func inputfileTransform() error {
	var err error
	var poCustomers = biz.PoCustomerList
	transf := poCustomers.GetTransformFunc(Inputtpl)
	if transf == nil {
		return fmt.Errorf("找不到转换对应的转换函数")
	}

	filesplit := strings.Split(Inputfile, ".")
	fileext := filesplit[len(filesplit)-1]
	outputfile := strings.Replace(Inputfile, "."+fileext, "-Done."+fileext, 1)
	fmt.Println("输入文件：", Inputfile, "输出文件：", outputfile)
	_, err = transf(Inputtpl, Inputfile, outputfile)
	if err != nil {
		fmt.Println("转换失败:", err)
	}
	fp, _ := filepath.Abs(Inputfile)
	exec.Command("cmd", "/c", "start", filepath.Dir(fp)).Start()
	return err
}
