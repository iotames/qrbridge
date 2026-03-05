package webserver

import (
	"fmt"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"

	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/webserver/amis"
)

func getAmisPoImportPage(ctx httpsvr.Context) {
	pageConf := amis.NewPage(getAmisPageTitle(ctx.Request.URL.Path))
	item1 := amis.NewFormItem().Set("label", "客户简称").Set("type", "select").Set("name", "inputtpl").Set("value", poCustomers[0].Code).Set("source", "/api/customer/list")
	item2 := amis.NewFormItem().Set("type", "input-file").Set("name", "inputfile").Set("accept", ".xlsx").Set("label", "上传.xlsx文件").Set("maxSize", 10048576).Set("receiver", "/api/uploadfile")
	pageConf.Body = *amis.NewForm("/api/potransform").AddItem(item1).AddItem(item2)
	ctx.Writer.Write(response.NewApiData(pageConf.Map(), "success", 0).Bytes())
}

func getAmisCmdConfig(ctx httpsvr.Context) {
	var ok bool
	title := getAmisPageTitle(ctx.Request.URL.Path)

	do := ctx.GetQueryValue("do", "")
	if do != "" {
		domap := map[string]string{
			"sync": "数据同步",
		}
		if title, ok = domap[do]; !ok {
			title = "快捷操作"
		}

		item1 := amis.NewFormItem().Set("label", "同步类型").Set("type", "select").
			Set("name", "optname").Set("value", "userlist").AddSelectOption("人员同步", "userlist").AddSelectOption("测试", "debug")
		// .Set("source", "/api/customer/list")
		form1 := *amis.NewForm("/api/cmd/exec").AddItem(item1)

		grid1 := amis.NewGrid()
		grid1.Col(form1, 3)

		// "ws://localhost:8777"
		wsaddr := fmt.Sprintf("ws://127.0.0.1:%d", conf.WebSocketPort)
		// customComp  := amis.NewWebSocket(wsaddr)
		customComp := amis.BuildWebSocketCustom(wsaddr, "cmd-output-area")
		grid1.Col(customComp.Map(), 9)

		// grid2 := amis.NewGrid()
		// grid2.Col(customComp.Map(), 9)

		page := amis.NewPage(title)
		// grids := []*amis.Grid{grid1, grid2}
		grids := []*amis.Grid{grid1}
		page.AddBody(grids)
		ctx.Writer.Write(response.NewApiData(page.Map(), "success", 0).Bytes())
	} else {
		ctx.Json(map[string]any{"code": 400, "msg": "do参数不能为空", "title": title}, 200)
	}

}
