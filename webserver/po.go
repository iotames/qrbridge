package webserver

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/service"
)

func poimport(ctx httpsvr.Context) {
	var err error
	var b2 []byte
	// var b1 io.ReadCloser
	// b1, err = ctx.Request.Body() runtime error: invalid memory address or nil pointer dereference

	// 读取请求体中的数据
	b2, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError("读取请求体失败io.ReadAll: " + err.Error()).Bytes())
		return
	}
	// 解析JSON数据
	var requestData map[string]interface{}
	err = json.Unmarshal(b2, &requestData)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError("解析JSON失败: " + err.Error()).Bytes())
		return
	}

	// 获取filepath字段
	filepath, ok := requestData["filepath"].(string)
	if !ok || filepath == "" {
		ctx.Writer.Write(response.NewApiDataServerError("filepath字段错误").Bytes())
		return
	}

	// 打印filepath字段
	fmt.Printf("接收到的filepath: %s\n", filepath)
	f, err := service.NewTableFile(filepath).OpenExcel()
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError("打开Excel文件失败: " + err.Error()).Bytes())
		return
	}
	sheets := f.GetSheetList()
	for i, sheet := range sheets {
		// 第57行为标题行
		// 客户款号C58
		// 款式标题G58 Kay High Waist Tights  Black/Black, XS,  0101
		noTitle, _ := f.GetCellValue(sheet, "C57")
		no, _ := f.GetCellValue(sheet, "C58")
		fmt.Printf("---sheet(%d-%s)----no(%s - %s)-------\n", i, sheet, noTitle, no)
	}
	defer f.Close()
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"filepath": filepath}, "success", 200).Bytes())
}
