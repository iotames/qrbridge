package webserver

import (
	"net/http"
	"strings"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/conf"
)

func setMiddlewares(svr *httpsvr.EasyServer) {
	svr.AddMiddleHead(NewMiddleCORS("*"))
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

// middleCORS CORS跨域设置中间件
type middleCORS struct {
	allowOrigin string
}

// NewMiddleCORS CORS中间件: 跨域设置。例: NewMiddleCORS("*")
// allowOrigin: 允许跨域的站点。默认值为 "*"。可将将 * 替换为指定的域名
func NewMiddleCORS(allowOrigin string) *middleCORS {
	if allowOrigin == "" {
		allowOrigin = "*"
	}
	return &middleCORS{allowOrigin: allowOrigin}
}

func (m middleCORS) Handler(w http.ResponseWriter, r *http.Request, dataFlow *httpsvr.DataFlow) (subNext bool) {
	w.Header().Set("Access-Control-Allow-Origin", m.allowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Accept, Token, Auth-Token, X-Requested-With")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Private-Network", "true")
	// 注意：如果需要携带Cookie等凭证，则不能使用 Access-Control-Allow-Origin: *，必须指定具体来源，并设置 Access-Control-Allow-Credentials: true。
	// 根据W3C CORS标准，当 Access-Control-Allow-Credentials为 true时，Access-Control-Allow-Origin绝对不能为 *，必须是一个明确的、单一的源。
	dataFlow.SetDataReadonly("CorsAllowOrigin", m.allowOrigin)
	return r.Method != "OPTIONS"
}
