package webserver

import (
	"fmt"

	"github.com/iotames/easyserver/httpsvr"
)

var webServerPort int

func Run(port int) {
	svr := httpsvr.NewEasyServer(fmt.Sprintf(":%d", port))
	webServerPort = port
	setMiddlewares(svr)
	setHandler(svr)
	svr.ListenAndServe()
}
