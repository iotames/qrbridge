package webserver

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

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

func poimport(ctx httpsvr.Context) {
	inputtpl, inputfile, outputfile, err := poimportCheck(ctx)
	if err != nil {
		return
	}
	// 打印inputfile字段
	fmt.Printf("接收到的inputfile(%s); outputfile(%s)\n", inputfile, outputfile)
	tpllist := []string{"Rohnisch"}
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
	outputfile := strings.Replace(inputfile, ".xlsx", "-Done.xlsx", 1)

	// 打印inputfile字段
	fmt.Printf("接收到的inputfile(%s); outputfile(%s)\n", inputfile, outputfile)
	tpllist := []string{"Rohnisch"}
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

func uploadfile(ctx httpsvr.Context) {
	// 获取上传的文件
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// fmt.Printf("----------upload-----filename(%s)----\n", header.Filename)
	defer file.Close()
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

	// 构建完整文件路径并保存文件
	filepath := filepath.Join(uploadDir, saveFilename)
	dst, err := os.Create(filepath)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError("创建目标文件失败: " + err.Error()).Bytes())
		return
	}
	defer dst.Close()

	// 将上传文件内容复制到目标文件
	if _, err := io.Copy(dst, file); err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError("保存文件失败: " + err.Error()).Bytes())
		return
	}
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"value": filepath}, "success", 0).Bytes())
	// {
	//   "status": 0,
	//   "msg": "",
	//   "data": {
	//     "value": "xxxx"
	//   }
	// }
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// fmt.Println(stat.IsDir())
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
