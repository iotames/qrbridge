package webserver

import (
	"net/http"
	"strings"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/conf"
)

func setMiddlewares(svr *httpsvr.EasyServer) {
	svr.AddMiddleHead(httpsvr.NewMiddleCORS("*"))
	svr.AddMiddleHead(httpsvr.NewMiddleStatic("/static/amis", "./resource/amis"))
	svr.AddMiddleHead(UserAuthMiddle{})
}

type UserAuthMiddle struct{}

// 自定义用户中间件：进行用户认证
// ‌401 Unauthorized‌：表示‌未提供身份凭证‌或凭证无效，需先登录。
// ‌403 Forbidden‌：表示‌已登录但无权限‌，即使提供正确凭证仍被拒绝。
func (h UserAuthMiddle) Handler(w http.ResponseWriter, r *http.Request, dataFlow *httpsvr.DataFlow) (next bool) {
	if strings.HasPrefix(r.URL.Path, "/user/") {
		token := r.URL.Query().Get("token")
		if token == "" {
			w.Write(response.NewApiDataFail("权限令牌不能为空", 401).Bytes())
			return false
		}
		if conf.AppToken == "" {
			w.Write(response.NewApiDataFail("权限令牌APP_TOKEN未配置", 500).Bytes())
			return false
		}
		if conf.AppToken == token {
			return true
		}
		w.Write(response.NewApiDataFail("权限错误", 401).Bytes())
		return false
	}
	return true
}
