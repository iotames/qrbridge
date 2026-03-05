package amis

type WsConf struct {
	Url  string      `json:"url"`
	Data JsonContent `json:"data"`
}

type WebSocket struct {
	BaseComponent
	Ws   WsConf      `json:"ws"`
	Body JsonContent `json:"body"`
}

func NewWebSocket(url string) *WebSocket {
	ws := &WebSocket{BaseComponent: BaseComponent{Type: "service"}}
	ws.Ws.Url = url
	return ws
}

func (w WebSocket) Map() map[string]any {
	return map[string]any{
		"type": w.Type,
		"ws":   w.Ws,
		"body": w.Body,
	}
}
