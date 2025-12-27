package biz

import (
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA5ygcTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, 0, poSheetDataParseA5ygc)
}

// 从Excel的每个sheet页面解析数据
func poSheetDataParseA5ygc(f *excelize.File, sheetName string, info *PoInfo) error {
	ids := getOkOrderItemRowIndexs(f, sheetName, "E", 2, 3, []string{"SKU"})
	sizeMap := map[string]string{
		"00": "XXS",
		"01": "XS",
		"02": "S",
		"03": "M",
		"04": "L",
		"05": "XL",
		"06": "XXL",
		"07": "XXXL",
	}
	for _, rowindex := range ids {
		item := OrderItem{}
		sku := getCellTrimSpace(f, sheetName, "E", rowindex) // sku ok
		if len(sku) > 10 {
			item.StyleNo = sku[0:12]    // 客户款号 ok
			sizeKey := sku[len(sku)-2:] // 取E列的后两位转换成尺码 00：XXS 01：XS …… 07：XXXL
			item.Size = sizeMap[sizeKey]
		}
		item.Desc = getCellTrimSpace(f, sheetName, "D", rowindex) // ok Soft-Motion 7/8 Legging in Black
		item.ColorEn = getA5ygcColorByDesc(item.Desc)             // 获取颜色尺码 ok
		item.PoNo = getCellTrimSpace(f, sheetName, "G", rowindex) // ok 客人没有PO的要求，直接取G列的合同号作为PO NO

		// 离厂交期 客户交期 H列 15/01/2026
		deliveryDateCustomerTxt := getCellTrimSpace(f, sheetName, "H", rowindex)       // 获取客户交期。15/01/2026
		deliveryDateCustomer, err := time.Parse("02/01/2006", deliveryDateCustomerTxt) // 客户交期。必填。
		// 字符串格式：2006-01-02。默认值为空字符串
		deliveryDateCustomerStr := ""     // 客户交期。必填。
		deliveryDateFactoryLeaveStr := "" // 离厂交期
		deliveryDateFactoryStr := ""      // 工厂交期

		if err == nil {
			// deliveryDateFactoryLeave := deliveryDateCustomer.AddDate(0, 0, -7)          // 离厂交期。必填。客户交期-7天
			deliveryDateFactoryLeave := deliveryDateCustomer                            // 离厂交期=客户交期
			deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7)           // 工厂交期。非必填。离厂交期-7天
			deliveryDateCustomerStr = deliveryDateCustomer.Format("2006-01-02")         // 客户交期
			deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02") // 离厂交期
			deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")           // 工厂交期
		}
		// 特殊处理，离厂交期和客户交期一样。
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。客户交期-7天
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天

		// TODO 1、取F列数据 2、要求将同SKU(E列)合并数量，出货时再根据客户要求手动拆分
		qtyStr := getCellTrimSpace(f, sheetName, "F", rowindex)          // ok 原始数据。订单数量。必填
		fmt.Sscanf(qtyStr, "%d", &item.Qty)                              // ok 转换为整型。订单数量。必填
		item.DestCountry = getCellTrimSpace(f, sheetName, "A", rowindex) // ok 目的国。必填。
		info.OrderItems = append(info.OrderItems, item)
		fmt.Printf("----sheet(%s)---rowindex(%d)---orderItem(%+v)------\n", sheetName, rowindex, item)
	}
	return nil
}

func getA5ygcColorByDesc(desc string) string {
	descsplit := strings.Split(desc, "in")
	splitlen := len(descsplit)
	if splitlen > 1 {
		return strings.TrimSpace(descsplit[splitlen-1])
	}
	return ""
}
