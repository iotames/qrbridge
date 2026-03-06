package webserver

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/tcpserver"
)

var cmdLastExeAt time.Time

// /user/cmd?do=sync&token=xxxxxxxx
func execmd(ctx httpsvr.Context) {
	var err error
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
		// 5分钟内不要重复提交
		if optname == "userlist" {
			if time.Since(cmdLastExeAt) < time.Minute {
				// http状态码设置为400也可以工作 ctx.Json(map[string]any{"status": http.StatusBadRequest, "msg": "请求已提交, 5分钟后再试"}, 400)
				ctx.Json(map[string]any{"status": http.StatusBadRequest, "msg": "请求已提交, 5分钟后再试"}, 200)
				return
			}
			cmdLastExeAt = time.Now()
		}
		err = execByName(optname)
		if err != nil {
			log.Printf("---execmd--postdata(%+v)-error(%+v)----\n", postdata, err)
			ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
			return
		}
	} else {
		ctx.Writer.Write(response.NewApiDataFail("参数错误", 400).Bytes())
		return
	}
	ctx.Writer.Write(response.NewApiDataOk("执行成功").Bytes())
}

func execByName(optname string) error {
	var cmd *exec.Cmd
	var connWriters []io.Writer
	wsvr := tcpserver.GetServer()
	if wsvr != nil {
		connWriters = wsvr.GetOutputWriters()
	}

	switch optname {
	case "userlist":
		cmd = exec.Command("/bin/bash", "-c", "/home/santic/kettle_hour.sh")
	case "debug":
		timeText := time.Now().Format("2006-01-02 15:04:05")
		if runtime.GOOS == "windows" {
			// // 二进制帧
			// cmd = exec.Command("cmd", "/c", "ping baidu.com")
			// 文本帧
			cmd = exec.Command("cmd", "/c", "echo Hello Santic "+timeText)
		} else {
			// cmd = exec.Command("/bin/bash", "-c", "ping baidu.com")
			// 文本帧
			cmd = exec.Command("/bin/bash", "-c", "echo Hello Santic "+timeText)
		}
	}
	// 标准输出：同时写入本地 stdout 和所有连接
	stdoutWriters := append([]io.Writer{os.Stdout}, connWriters...)
	cmd.Stdout = io.MultiWriter(stdoutWriters...)

	// 标准错误：同时写入本地 stderr 和所有连接
	stderrWriters := append([]io.Writer{os.Stderr}, connWriters...)
	cmd.Stderr = io.MultiWriter(stderrWriters...)

	// 启动命令
	err := cmd.Start()
	if err != nil {
		return err
	}

	// 关键修复：启动一个协程来等待命令结束，以回收系统资源，避免僵尸进程
	go func() {
		waitErr := cmd.Wait()
		if waitErr != nil {
			// 此处可以记录日志，但不应影响主流程返回
			log.Printf("---execByName--命令等待结束时报错(%v)---\n", waitErr)
		}
	}()

	return nil // 主函数立即返回，不阻塞HTTP请求
}
