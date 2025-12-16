package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/biz"
)

var PoCustomers = []biz.PoCustomer{
	{Code: "A89SP", Remark: "Rohnisch"},
	{Code: "A5YGC", Remark: "A5YGC"},
	// {"A6WHM", "A6WHM"},
	// {"B1ZTV", "B1ZTV"},
	// {"AH8SW", "AH8SW"},
	// {"A63AM", "A63AM"},
}

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
	options := make([]map[string]string, len(PoCustomers))
	for i, v := range PoCustomers {
		options[i] = map[string]string{"label": v.Code, "value": v.Code}
	}
	// json返回
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"options": options}, "success", 0).Bytes())
}
