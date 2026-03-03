package webserver

import (
	"os"
	"os/exec"
	"time"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

var cmdLastExeAt time.Time

// /user/cmd?do=sync&token=ioqwuyhfkluhdsflplqxzbvjhn
func execmd(ctx httpsvr.Context) {
	var err error
	// 1分钟内不要重复提交
	if time.Since(cmdLastExeAt) < time.Minute {
		// ctx.Writer.Write(response.NewApiDataFail(, 400).Bytes())
		ctx.Json(map[string]any{"status": 404, "msg": "请求已提交, 5分钟后再试"}, 200)
		return
	}
	cmdLastExeAt = time.Now()

	// postdata := map[string]string{}
	// err = ctx.GetPostJson(postdata)
	// 解析JSON失败json.Unmarshal error: json: Unmarshal(non-pointer map[string]string)
	var postdata map[string]string
	err = ctx.GetPostJson(&postdata)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataFail(err.Error(), 500).Bytes())
		return
	}
	optname, ok := postdata["optname"]
	if ok {
		err = execByName(optname)
		if err != nil {
			ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
			return
		}
	}
	ctx.Writer.Write(response.NewApiDataOk("执行成功").Bytes())
}

func execByName(optname string) error {
	var cmd *exec.Cmd
	switch optname {
	case "userlist":
		cmd = exec.Command("/bin/bash", "-c", "/home/santic/kettle_hour.sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Start()
}
