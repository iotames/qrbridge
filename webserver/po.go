package webserver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/miniutils"
)

// poimport 通过HTTP的API接口接收POST方法请求的JSON数据。包含：inputtpl, inputfile, outputfile三个字段。
// 接口返回数据：
//
//	{"code":200,"msg":"success","data":{"inputfile": inputfile, "outputfile": outputfile}}
func poimport(ctx httpsvr.Context) {
	requestData, err := checkJsonField(ctx, "inputtpl", "inputfile", "outputfile")
	if err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	inputtpl := requestData["inputtpl"].(string)
	inputfile := requestData["inputfile"].(string)
	outputfile := requestData["outputfile"].(string)

	// 打印inputfile字段
	fmt.Printf("接收到的inputfile(%s); outputfile(%s)\n", inputfile, outputfile)
	tpllist := poCustomers.GetCodeList()
	if !slices.Contains(tpllist, inputtpl) {
		err = fmt.Errorf("inputtpl参数错误: 仅支持(%s)", strings.Join(tpllist, ","))
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	potransformfunc := poCustomers.GetTransformFunc(inputtpl)
	if potransformfunc == nil {
		err = fmt.Errorf("%s 找不到转换对应的转换函数", inputfile)
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	var tempInputfile string
	if strings.HasPrefix(inputfile, "http") {
		unixtime := time.Now().Unix()
		saveFilename := fmt.Sprintf("%s-%d.xlsx", time.Now().Format(time.DateOnly), unixtime)

		// 创建上传目录（如果不存在）
		uploadDir := "runtime/upload"
		if !IsPathExists(uploadDir) {
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				ctx.Writer.Write(response.NewApiDataServerError("创建上传目录失败：" + err.Error()).Bytes())
				return
			}
		}
		// 定义完整的临时文件路径
		tempInputfile = filepath.Join(uploadDir, saveFilename)
		err = miniutils.NewHttpRequest(inputfile).Download(tempInputfile)
		if err != nil {
			ctx.Writer.Write(response.NewApiDataServerError("下载文件失败：" + err.Error()).Bytes())
			return
		}
		inputfile = tempInputfile
	}

	_, err = potransformfunc(inputfile, outputfile)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	if miniutils.IsPathExists(tempInputfile) {
		fmt.Println("删除临时文件:", tempInputfile)
		err = os.Remove(tempInputfile)
		if err != nil {
			fmt.Println("删除临时文件失败:", err)
		} else {
			fmt.Println("删除临时文件成功:", tempInputfile)
		}
	} else {
		fmt.Println("tempInputfile", tempInputfile, "不存在")
	}

	ctx.Writer.Write(response.NewApiData(response.JsonObject{"inputfile": inputfile, "outputfile": outputfile}, "success", 200).Bytes())
}

// potransform 通过HTTP的API接口接收POST方法请求的JSON数据。包含：inputtpl, inputfile两个字段。
func potransform(ctx httpsvr.Context) {
	requestData, err := checkJsonField(ctx, "inputtpl", "inputfile")
	if err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	inputtpl := requestData["inputtpl"].(string)
	inputfile := requestData["inputfile"].(string)

	// 获取outputfile字段
	outputfileBase := filepath.Base(inputfile)
	outputfileBaseNew := inputtpl + "-" + strings.Replace(outputfileBase, ".xlsx", "-Done.xlsx", 1)
	outputfile := strings.Replace(inputfile, outputfileBase, outputfileBaseNew, 1)

	// 打印inputfile字段
	fmt.Printf("接收到的inputfile(%s); outputfile(%s)\n", inputfile, outputfile)
	tpllist := poCustomers.GetCodeList()
	if !slices.Contains(tpllist, inputtpl) {
		err = fmt.Errorf("inputtpl参数错误: 仅支持(%s)", strings.Join(tpllist, ","))
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	fmt.Println("outputfile", outputfile)

	potransformfunc := poCustomers.GetTransformFunc(inputtpl)
	if potransformfunc == nil {
		err = fmt.Errorf("%s 找不到转换对应的转换函数", inputfile)
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// miniutils.NewHttpRequest("https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png").Download("runtime/baidu.png")
	_, err = potransformfunc(inputfile, outputfile)

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
