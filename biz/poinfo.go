package biz

import (
	"time"
)

// A89SP 谢小玉 Rohnisch
// A5YGC 张春梅
// A6WHM 张真真
// B1ZTV 王芳
// AH8SW 李丝婷
// A63AM 陈雪娇
// BEWCW 林容情
// A8ASO	王阿玲

type PoCustomer struct {
	Code            string
	Remark          string
	poTransformFunc func(inputfile, outputfile string) (info PoInfo, err error)
}

type PoInfo struct {
	PoNo, DestCountry, DestPortName string
	DeliveryDateCustomer            time.Time
	OrderItems                      []OrderItem
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
	{Code: "BEWCW", Remark: "Icaniwill|ICIW", poTransformFunc: PoBewcwTransform},
	{Code: "A6WHM", Remark: "HEMA", poTransformFunc: PoA6whmTransform},
	{Code: "B1ZTV", Remark: "TEVEO", poTransformFunc: PoB1ztvTransform},
	{Code: "AH8SW", Remark: "ALPHA", poTransformFunc: PoAh8swTransform},
	{Code: "A63AM", Remark: "A63AM", poTransformFunc: PoA63amTransform},
	{Code: "A6TON", Remark: "A6TON", poTransformFunc: PoA6tonTransform},
	{Code: "A8ASO", Remark: "STRONGER", poTransformFunc: PoA8asoTransform},
}

func (pc PoCustomers) GetCodeList() []string {
	var list []string
	for _, v := range pc {
		list = append(list, v.Code)
	}
	return list
}

// GetTransformFunc 根据code获取转换函数
// code为inputtpl转换模板名称，也是客户简称
func (pc PoCustomers) GetTransformFunc(code string) func(inputfile, outputfile string) (info PoInfo, err error) {
	for _, v := range pc {
		if v.Code == code {
			return v.poTransformFunc
		}
	}
	return nil
}
