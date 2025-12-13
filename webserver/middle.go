package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
)

func setMiddlewares(svr *httpsvr.EasyServer) {
	svr.AddMiddleHead(httpsvr.NewMiddleCORS("*"))
	svr.AddMiddleHead(httpsvr.NewMiddleStatic("/static/amis", "./resource/amis"))
}
