package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
)

func Run(addr string) {
	svr := httpsvr.NewEasyServer(addr)
	setMiddlewares(svr)
	setHandler(svr)
	svr.ListenAndServe()
}
