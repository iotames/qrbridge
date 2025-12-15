package biz

// A89SP 谢小玉 Rohnisch
// A5YGC 张春梅
// A6WHM 张真真
// B1ZTV 王芳
// AH8SW 李丝婷
// A63AM 陈雪娇

type PoCustomer struct {
	Code   string
	Remark string
}

var PoCustomers = []PoCustomer{
	{"A89SP", "Rohnisch"},
	{"A5YGC", "A5YGC"},
	// {"A6WHM", "A6WHM"},
	// {"B1ZTV", "B1ZTV"},
	// {"AH8SW", "AH8SW"},
	// {"A63AM", "A63AM"},
}

type PoInfo struct {
	OrderItems []OrderItem
}

type OrderItem struct {
	StyleNo                  string // 必填 款号
	Color                    string // 必填 颜色
	ColorEn                  string // 必填 英文颜色
	ColorNo                  string // 色号
	PoNo                     string // 必填
	Size                     string // 必填
	DeliveryDateFactory      string // 工厂交期
	DeliveryDateFactoryLeave string // 必填 离场交期
	DeliveryDateCustomer     string // 必填 客户交期
	Qty                      int    // 必填 订单数量
	DestCountry              string // 必填	目的国
	DestPortName             string // 目的港
	Desc                     string
}
