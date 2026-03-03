package amis

type FormItem map[string]any
type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func NewFormItem() *FormItem {
	item := make(FormItem)
	return &item
}

func (f *FormItem) Set(k string, v any) *FormItem {
	(*f)[k] = v // 需要解引用
	return f
}

func (f *FormItem) AddSelectOption(label, value string) *FormItem {
	opt := SelectOption{label, value}
	opts, ok := (*f)["options"]
	if !ok {
		opts = []SelectOption{opt}
	} else {
		opts = append(opts.([]SelectOption), opt)
	}
	(*f)["options"] = opts
	return f
}

// item1 := amis.NewFormItem().Set("label", "客户简称").Set("type", "select").Set("name", "inputtpl").Set("value", PoCustomers[0].Code).Set("source", "/api/customer/list")
// item2 := amis.NewFormItem().Set("type", "input-file").Set("name", "inputfile").Set("accept", ".xlsx").Set("label", "上传.xlsx文件").Set("maxSize", 10048576).Set("receiver", "/api/uploadfile")
