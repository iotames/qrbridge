package biz

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
