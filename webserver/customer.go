package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/biz"
)

var poCustomers = biz.PoCustomerList

// 客户列表
//
//	{
//	"status": 0,
//	"data": {
//		"options": [
//		 {"label": "A89SP", "value": "A89SP"},
//		 {"label": "A5YGC", "value": "A5YGC"}
//		]
//	 }
//	}
func customerList(ctx httpsvr.Context) {
	// 获取所有客户名称
	options := make([]map[string]string, len(poCustomers))
	for i, v := range poCustomers {
		options[i] = map[string]string{"label": v.Code, "value": v.Code}
	}
	// json返回
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"options": options}, "success", 0).Bytes())
}
