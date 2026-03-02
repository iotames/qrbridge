package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

func execmd(ctx httpsvr.Context) {
	// // http://172.16.160.11/user/cmd?do=sync&token=ioqwuyhfkluhdsflplqxzbvjhn
	// /home/santic/kettle_hour.sh

	ctx.Writer.Write(response.NewApiDataOk("执行成功").Bytes())
}
