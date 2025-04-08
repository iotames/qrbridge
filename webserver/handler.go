package webserver

import (
	"encoding/json"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
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
	code := ctx.Request.URL.Query().Get("code")
	requestIp := util.GetHttpClientIP(ctx.Request)
	userAgent := ctx.Request.Header.Get("User-Agent")

	// 将请求头转换为JSON字符串
	requestHeaders, err := json.Marshal(ctx.Request.Header)
	if err != nil {
		lg.Errorf("convert headers to json failed: %v", err)
		requestHeaders = []byte("{}")
	}
	lg.Debugf("code: %s, request_ip: %s, user_agent(%s)---hdr(%s)", code, requestIp, userAgent, string(requestHeaders))
	ctx.Writer.Write(response.NewApiDataOk("hello your ip is: " + requestIp).Bytes())
}
