package conf

import (
	"fmt"
	"os"

	"github.com/iotames/easyconf"
	"github.com/iotames/miniutils"
)

var cf *easyconf.Conf

const DRIVER_MYSQL = "mysql"
const DRIVER_POSTGRES = "postgres"

const DEFAULT_ENV_FILE = ".env"
const DEFAULT_RUNTIME_DIR = "runtime"
const DEFAULT_RESOURCE_DIR = "resource"
const DEFAULT_WEB_SERVER_PORT = 8080
const DEFAULT_DB_DRIVER = DRIVER_POSTGRES
const DEFAULT_DB_HOST = "127.0.0.1"
const DEFAULT_DB_PORT = 5432
const DEFAULT_DB_NAME = "postgres"
const DEFAULT_DB_SCHEMA = "public"
const DEFAULT_DB_USERNAME = "postgres"
const DEFAULT_DB_PASSWORD = "postgres"

var RuntimeDir string
var ResourceDir string
var ToBaseUrl string
var WebServerPort int
var ShowSql bool
var DbDriver, DbHost, DbName, DbSchema, DbUsername, DbPassword string
var DbPort, DbNodeId, EncryptMultiple, EncryptAdd int

func getEnvFile() string {
	efile := os.Getenv("QR_ENV_FILE")
	if efile == "" {
		efile = DEFAULT_ENV_FILE
	}
	return efile
}

func setConfByEnv() {
	// # 设置 QR_ENV_FILE 环境变量，可更改配置文件路径。
	cf = easyconf.NewConf(getEnvFile())

	cf.StringVar(&RuntimeDir, "RUNTIME_DIR", DEFAULT_RUNTIME_DIR, "")
	cf.StringVar(&ResourceDir, "RESOURCE_DIR", DEFAULT_RESOURCE_DIR, "")
	cf.IntVar(&WebServerPort, "WEB_SERVER_PORT", DEFAULT_WEB_SERVER_PORT, "启动Web服务器的端口号")

	cf.BoolVar(&ShowSql, "SHOW_SQL", false, "是否输出SQL调试信息")
	cf.StringVar(&DbDriver, "DB_DRIVER", DEFAULT_DB_DRIVER, "数据库类型: mysql,sqlite3,postgres")
	cf.StringVar(&DbHost, "DB_HOST", DEFAULT_DB_HOST, "数据库主机地址")
	cf.StringVar(&DbName, "DB_NAME", DEFAULT_DB_NAME, "数据库名")
	cf.StringVar(&DbSchema, "DB_SCHEMA", DEFAULT_DB_SCHEMA, "数据库schema")
	cf.IntVar(&DbPort, "DB_PORT", DEFAULT_DB_PORT, "数据库端口号:5432(postgres);3306(mysql)")
	cf.StringVar(&DbUsername, "DB_USERNAME", DEFAULT_DB_USERNAME, "数据库用户名")
	cf.StringVar(&DbPassword, "DB_PASSWORD", DEFAULT_DB_PASSWORD, "数据库密码")
	cf.IntVar(&DbNodeId, "DB_NODE_ID", 1, "数据库节点号")

	cf.IntVar(&EncryptMultiple, "ENCRYPT_MULTIPLE", 2, "加密倍数")
	cf.IntVar(&EncryptAdd, "ENCRYPT_ADD", 10086, "加密增量")
	cf.StringVar(&ToBaseUrl, "TO_BASE_URL", "", "跳转目标的URL前缀。如：https://www.baidu.com")

	cf.Parse()
}

func LoadEnv() error {
	var err error
	setConfByEnv()
	if !miniutils.IsPathExists(ResourceDir) {
		err = miniutils.Mkdir(ResourceDir)
		if err != nil {
			return err
		}
	}
	if !miniutils.IsPathExists(RuntimeDir) {
		fmt.Printf("------创建runtime目录(%s)--\n", RuntimeDir)
		err = os.Mkdir(RuntimeDir, 0755)
		if err != nil {
			fmt.Printf("----runtime目录(%s)创建失败(%v)---\n", RuntimeDir, err)
		}
	}
	return err
}
