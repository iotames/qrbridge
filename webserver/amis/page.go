package amis

// PageConfig 定义 Amis 页面配置结构
type PageConfig struct {
	Type  string      `json:"type"`
	Title string      `json:"title"`
	Body  JsonContent `json:"body"`
}

func NewPage(title string) *PageConfig {
	return &PageConfig{
		Type:  "page",
		Title: title,
	}
}

func (p *PageConfig) AddBody(body JsonContent) *PageConfig {
	p.Body = body
	return p
}

func (p PageConfig) Map() map[string]any {
	return map[string]any{
		"type":  p.Type,
		"title": p.Title,
		"body":  p.Body,
	}
}
