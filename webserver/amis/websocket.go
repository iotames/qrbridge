package amis

type WsConf struct {
	Url         string      `json:"url"`
	Data        JsonContent `json:"data"`
	ResponseKey string      `json:"responseKey"`
}

type WebSocket struct {
	BaseComponent
	Ws   WsConf      `json:"ws"`
	Body JsonContent `json:"body"`
}

func NewWebSocket(url string) *WebSocket {
	ws := &WebSocket{BaseComponent: BaseComponent{Type: "service"}}
	ws.Ws.Url = url
	// ws.Ws.Data = map[string]string{"name": "${name}"}
	ws.Ws.ResponseKey = "output" // 假设返回的字符串放在output字段
	ws.Body = map[string]any{
		"type": "tpl",
		"tpl":  "<pre style='background:#f5f5f5;padding:10px;'>${output}</pre>",
		// "type":       "log",
		// "height":     300, // 设置固定高度
		// "autoScroll": true,
		// "placement":  "top",
		// "source":     "${output}", // 绑定数据
	}
	return ws
}

func (w WebSocket) Map() map[string]any {
	return map[string]any{
		"type": w.Type,
		"ws":   w.Ws,
		"body": w.Body,
	}
}
