package webserver

import (
	"fmt"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/biz"
	"github.com/iotames/qrbridge/service"
)

func poimport(ctx httpsvr.Context) {
	var err error
	// 解析JSON数据
	var requestData map[string]string
	err = postJsonValue(ctx, &requestData)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	// 获取filepath字段
	filepath, ok := requestData["filepath"]
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
		biz.PoSheetRohnisch(f, sheet, i)
		if i > 3 {
			break
		}
	}
	defer f.Close()
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"filepath": filepath}, "success", 200).Bytes())
}
