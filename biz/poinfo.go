package biz

var poCommonTplTitleRow = []interface{}{"客户款号*", "颜色*", "英文颜色*", "色号", "PO NO*", "尺码*", "工厂交期", "离厂交期*", "客户交期*", "订单数量*", "目的国*", "目的港"}

// A89SP 谢小玉 Rohnisch
// A5YGC 张春梅
// A6WHM 张真真
// B1ZTV 王芳
// AH8SW 李丝婷
// A63AM 陈雪娇

type PoCustomer struct {
	Code            string
	Remark          string
	poTransformFunc func(inputtpl, inputfile, outputfile string) (info PoInfo, err error)
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

type PoCustomers []PoCustomer

var PoCustomerList = PoCustomers{
	{Code: "A89SP", Remark: "Rohnisch", poTransformFunc: PoRohnischTransform},
	{Code: "A5YGC", Remark: "A5YGC", poTransformFunc: PoA5ygcTransform},
	// {"A6WHM", "A6WHM"},
	// {"B1ZTV", "B1ZTV"},
	// {"AH8SW", "AH8SW"},
	// {"A63AM", "A63AM"},
}

func (pc PoCustomers) GetCodeList() []string {
	var list []string
	for _, v := range pc {
		list = append(list, v.Code)
	}
	return list
}

func (pc PoCustomers) GetTransformFunc(code string) func(inputtpl, inputfile, outputfile string) (info PoInfo, err error) {
	for _, v := range pc {
		if v.Code == code {
			return v.poTransformFunc
		}
	}
	return nil
}
