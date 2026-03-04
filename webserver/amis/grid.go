package amis

// https://baidu.github.io/amis/zh-CN/components/grid

type GridColumn struct {
	ClassName string        `json:"columnClassName"`
	Md        int           `json:"md"`
	Body      []JsonContent `json:"body"`
}

type Grid struct {
	// Type  string     `json:"type"`
	BaseComponent
	Columns []GridColumn `json:"columns"`
}

func NewGrid() *Grid {
	return &Grid{BaseComponent: BaseComponent{Type: "grid"}}
}

// Col 添加列
// md 宽度占比： 1 - 12
func (g *Grid) Col(body JsonContent, md int) *Grid {
	if body == nil {
		body = map[string]any{}
	}
	// { "md": 9, "body": [{...}]}
	g.Columns = append(g.Columns, GridColumn{Body: []JsonContent{body}, Md: md})
	return g
}

func (g Grid) Map() map[string]any {
	return map[string]any{
		"type":    g.Type,
		"columns": g.Columns,
	}
}
