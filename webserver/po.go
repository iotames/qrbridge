package webserver

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/biz"
)

func poimportCheck(ctx httpsvr.Context) (inputtpl string, inputfile string, outputfile string, err error) {
	var ok bool
	// 解析JSON数据
	var requestData map[string]string
	err = postJsonValue(ctx, &requestData)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	// 获取inputtpl字段
	inputtpl, ok = requestData["inputtpl"]
	if !ok || inputtpl == "" {
		// ctx.Writer.Write(response.NewApiDataFail("inputtpl字段错误", 400).Bytes())
		err = fmt.Errorf("inputtpl字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// 获取inputfile字段
	inputfile, ok = requestData["inputfile"]
	if !ok || inputfile == "" {
		err = fmt.Errorf("inputfile字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// 获取outputfile字段
	outputfile, ok = requestData["outputfile"]
	if !ok || outputfile == "" {
		err = fmt.Errorf("outputfile字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	return
}

// poimport 通过HTTP的API接口接收POST方法请求的JSON数据。包含：inputtpl, inputfile, outputfile三个字段。
// 接口返回数据：
//
//	{"code":200,"msg":"success","data":{"inputfile": inputfile, "outputfile": outputfile}}
func poimport(ctx httpsvr.Context) {
	inputtpl, inputfile, outputfile, err := poimportCheck(ctx)
	if err != nil {
		return
	}
	// 打印inputfile字段
	fmt.Printf("接收到的inputfile(%s); outputfile(%s)\n", inputfile, outputfile)
	tpllist := []string{"A89SP"}
	if !slices.Contains(tpllist, inputtpl) {
		err = fmt.Errorf("inputtpl参数错误: 仅支持(%s)", strings.Join(tpllist, ","))
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	_, err = biz.PoFileTransform(inputtpl, inputfile, outputfile)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"inputfile": inputfile, "outputfile": outputfile}, "success", 200).Bytes())
}

func potransform(ctx httpsvr.Context) {
	var ok bool
	var err error
	var inputtpl, inputfile string
	// 解析JSON数据
	var requestData map[string]string
	err = postJsonValue(ctx, &requestData)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	// 获取inputtpl字段
	inputtpl, ok = requestData["inputtpl"]
	if !ok || inputtpl == "" {
		// ctx.Writer.Write(response.NewApiDataFail("inputtpl字段错误", 400).Bytes())
		err = fmt.Errorf("inputtpl字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// 获取inputfile字段
	inputfile, ok = requestData["inputfile"]
	if !ok || inputfile == "" {
		err = fmt.Errorf("inputfile字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}

	// 获取outputfile字段
	outputfileBase := filepath.Base(inputfile)
	outputfileBaseNew := inputtpl + "-" + strings.Replace(outputfileBase, ".xlsx", "-Done.xlsx", 1)
	outputfile := strings.Replace(inputfile, outputfileBase, outputfileBaseNew, 1)

	// 打印inputfile字段
	fmt.Printf("接收到的inputfile(%s); outputfile(%s)\n", inputfile, outputfile)
	tpllist := []string{"A89SP"}
	if !slices.Contains(tpllist, inputtpl) {
		err = fmt.Errorf("inputtpl参数错误: 仅支持(%s)", strings.Join(tpllist, ","))
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	fmt.Println("outputfile", outputfile)
	_, err = biz.PoFileTransform(inputtpl, inputfile, outputfile)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}

	// 判断如果是Windows系统，则打开文件所在目录
	if runtime.GOOS == "windows" {
		// 获取输出文件的绝对路径，并提取其所在目录
		absPath, err := filepath.Abs(outputfile)
		if err != nil {
			// 处理获取绝对路径的错误，例如记录日志
			fmt.Printf("获取绝对路径失败: %v\n", err)
			return // 或继续执行，不中断主流程
		}
		dir := filepath.Dir(absPath)

		// 执行命令打开目录
		cmd := exec.Command("cmd", "/c", "start", dir)
		err = cmd.Start() // 使用Start而非Run，以便非阻塞地打开浏览器
		if err != nil {
			// 处理命令执行错误，例如记录日志
			fmt.Printf("打开目录失败: %v\n", err)
		}
		// 注意：这里使用Start，命令会在后台异步执行，我们通常不需要等待它结束
	}

	// json返回
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"inputfile": inputfile, "outputfile": outputfile}, "success", 0).Bytes())

	// // 文件下载逻辑
	// file, err := os.Open(outputfile)
	// if err != nil {
	// 	ctx.Writer.Write(response.NewApiDataServerError("无法打开生成的文件: " + err.Error()).Bytes())
	// 	return
	// }
	// defer file.Close()

	// // 获取文件信息
	// fileInfo, err := file.Stat()
	// if err != nil {
	// 	ctx.Writer.Write(response.NewApiDataServerError("无法获取文件信息: " + err.Error()).Bytes())
	// 	return
	// }
	// ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(outputfile))
	// ctx.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	// ctx.Writer.Header().Set("Last-Modified", fileInfo.ModTime().UTC().Format(time.RFC1123))
	// http.ServeContent(ctx.Writer, ctx.Request, filepath.Base(outputfile), fileInfo.ModTime(), file)
}
