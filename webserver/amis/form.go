package amis

type Form struct {
	Type  string `json:"type"`
	Title string `json:"title"`

	Mode string      `json:"mode"`
	Api  string      `json:"api"`
	Body []*FormItem `json:"body"`
}

func NewForm(apiurl, title string) *Form {
	return &Form{
		Type:  "form",
		Title: title,
		Mode:  "horizontal",
		Api:   apiurl,
		Body:  make([]*FormItem, 0),
	}
}

func (f *Form) AddItem(item *FormItem) *Form {
	f.Body = append(f.Body, item)
	return f
}

func (f Form) Map() map[string]any {
	return map[string]any{
		"type":  f.Type,
		"title": f.Title,
		"mode":  f.Mode,
		"api":   f.Api,
		"body":  f.Body,
	}
}
