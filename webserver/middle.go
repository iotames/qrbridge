package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
)

func setMiddlewares(svr *httpsvr.EasyServer) {
	svr.AddMiddleware(httpsvr.NewMiddleCORS("*"))
}
