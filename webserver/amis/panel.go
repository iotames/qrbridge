package amis

// https://baidu.github.io/amis/zh-CN/components/panel

type Panel struct {
	BaseComponent
	Title JsonContent `json:"title"`
	Body  JsonContent `json:"body"`
}

func NewPanel(title JsonContent, body JsonContent) *Panel {
	return &Panel{BaseComponent: BaseComponent{Type: "panel"}, Title: title, Body: body}
}

// Map 转换为 AMIS 配置
func (p *Panel) Map() map[string]any {
	return map[string]any{
		"type":  p.Type,
		"title": p.Title,
		"body":  p.Body,
	}
}

// {
//   "type": "page",
//   "body": {
//     "type": "panel",
//     "title": "面板标题",
//     "body": "面板内容"
//   }
// }
