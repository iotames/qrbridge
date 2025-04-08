package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/service"
	"github.com/iotames/qrbridge/util"
)

func setHandler(svr *httpsvr.EasyServer) {
	svr.AddHandler("GET", "/hello", hello)
	svr.AddHandler("GET", "/qrcode", qrcode)
}

func hello(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("hello api").Bytes())
}

func qrcode(ctx httpsvr.Context) {
	lg := util.GetLogger()
	// TODO
	service.GetOneQrcode(*ctx.Request)
	service.UpdateQrcode(*ctx.Request, true)
	lg.Debugf("qrcode-----")
	ctx.Writer.Write(response.NewApiDataOk("hello").Bytes())
}
