package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/webserver/amis"
)

func getAmisPageConfig(ctx httpsvr.Context) {
	pageConf := amis.NewPage("客户PO文件格式转换")
	item1 := amis.NewFormItem().Set("label", "客户简称").Set("type", "select").Set("name", "inputtpl").Set("value", poCustomers[0].Code).Set("source", "/api/customer/list")
	item2 := amis.NewFormItem().Set("type", "input-file").Set("name", "inputfile").Set("accept", ".xlsx").Set("label", "上传.xlsx文件").Set("maxSize", 10048576).Set("receiver", "/api/uploadfile")
	pageConf.Body = *amis.NewForm("/api/potransform").AddItem(item1).AddItem(item2)
	ctx.Writer.Write(response.NewApiData(pageConf.Json(), "success", 0).Bytes())
}
