package webserver

import (
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
	var err error
	title := getAmisPageTitle(ctx.Request.URL.Path)
	// var ok bool
	// do := ctx.GetQueryValue("do", "")
	// ctx.Json(map[string]any{"code": 400, "msg": "do参数不能为空", "title": title}, 200)
	// domap := map[string]string{
	// 	"sync": "数据同步",
	// }
	// if title, ok = domap[do]; !ok {
	// 	title = "快捷操作"
	// }

	item1 := amis.NewFormItem().Set("label", "操作类型").Set("type", "select").
		Set("name", "optname").Set("value", "userlist") // .Set("source", "/api/customer/list")
	cmds, err := GetCmds()
	if err != nil {
		ctx.Json(map[string]any{"code": 500, "msg": "获取命令列表错误：" + err.Error(), "title": title}, 500)
		return
	}
	for _, cmdinfo := range cmds {
		item1.AddSelectOption(cmdinfo.Title, cmdinfo.Name)
	}
	// .AddSelectOption("人员同步", "userlist").AddSelectOption("调试", "debug")

	form1 := *amis.NewForm("/api/cmd/exec").AddItem(item1)

	grid1 := amis.NewGrid()
	grid1.Col(form1, 3)
	// "ws://localhost:8777"
	wsaddr := conf.WebSocketAddr
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

}
