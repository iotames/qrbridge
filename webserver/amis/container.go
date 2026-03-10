package amis

// https://baidu.github.io/amis/zh-CN/components/container

type Container struct {
	BaseComponent
	Style map[string]any `json:"style,omitempty"` // 当字段值为‌零值时，该字段将被‌省略‌，不包含在输出的 JSON 中。
	Body  JsonContent
}

func NewContainer(body JsonContent) *Container {
	return &Container{
		BaseComponent: BaseComponent{Type: "container"},
		Body:          body,
	}
}

func (c *Container) StyleItem(k string, v any) *Container {
	if c.Style == nil {
		c.Style = make(map[string]any)
	}
	c.Style[k] = v
	return c
}

func (c *Container) Map() map[string]any {
	return map[string]any{
		"type":  c.Type,
		"style": c.Style,
		"body":  c.Body,
	}
}

// {
//   "type": "page",
//   "body": {
//     "type": "container",
//     "style": {
//       "backgroundColor": "#C4C4C4"
//     },
//     "body": "这里是容器内容区"
//   }
// }
