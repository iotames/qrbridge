package amis

type Form struct {
	Type string      `json:"type"`
	Mode string      `json:"mode"`
	Api  string      `json:"api"`
	Body []*FormItem `json:"body"`
}

func NewForm(apiurl string) *Form {
	return &Form{
		Type: "form",
		Mode: "horizontal",
		Api:  apiurl,
		Body: make([]*FormItem, 0),
	}
}

func (f *Form) AddItem(item *FormItem) *Form {
	f.Body = append(f.Body, item)
	return f
}
