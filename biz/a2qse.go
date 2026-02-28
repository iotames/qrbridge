package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA2qseTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, 0, PoSheetDataParseA2qse)
}

// 从Excel的单个sheet页面解析数据

// 颜色* 由业务员在标准模板中自己填写

// 色号 无

// 工厂交期 离厂交期-7天
// 离厂交期* 等于客户交期
// 客户交期* 取“XFTY Date”列下方的内容

// 目的国* 固定填USA
// 目的港 固定填Los Angeles

func PoSheetDataParseA2qse(f *excelize.File, sheetIndex int, info *PoInfo) error {

	var rows [][]string
	var err error
	var rowindex uint
	var qty int

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期

	info.DestCountry = "USA"
	info.DestPortName = "Los Angeles"

	sheetName := f.GetSheetName(sheetIndex)
	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}

	for i, row := range rows {
		// fmt.Printf("----PoSheetDataParseAh8sw--eachrow(%+v)---\n", row)
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// 跳出空数据行
		if strings.TrimSpace(row[0]) == "" {
			continue
		}
		// 跳过前6行
		if rowindex < 7 {
			continue
		}

		deliveryDateCustomerText := getCellTrimSpace(f, sheetName, "B", rowindex)
		// 1. 客户交期 2026/4/30, 04-30-26
		info.DeliveryDateCustomer, err = time.Parse("2006/1/02", deliveryDateCustomerText)
		if err != nil {
			// "02-Jan-06"对应：日-月缩写-年（2位）
			// "01-02-06" 对应：月-日-年（2位）
			info.DeliveryDateCustomer, _ = time.Parse("01-02-06", deliveryDateCustomerText)
			// TODO 解析更多时间格式
		}
		if !info.DeliveryDateCustomer.IsZero() {
			// 客户交期
			deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

			// 离厂交期
			deliveryDateFactoryLeave := info.DeliveryDateCustomer // 离厂交期=客户交期
			deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")

			// 工厂交期
			deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7) // 工厂交期=离厂交期-7天
			deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
		}
		// fmt.Printf("------deliveryDateCustomerText(%s)--deliveryDateCustomerStr(%s)---\n", deliveryDateCustomerText, info.DeliveryDateCustomer)

		// 2. PO NO* 取“PO Number”列下方的内容
		info.PoNo = getCellTrimSpace(f, sheetName, "A", rowindex)

		// 3. 客户款号*    "1、默认取SKU列最左侧的六位数字 2、如果SKU列为空，则取“New SKU”第一个“-”右侧的六位数字 3、开始前和陈银凤确认PO文件是否更新"
		styleNo := ""
		sku := getCellTrimSpace(f, sheetName, "G", rowindex)
		if len(sku) > 5 {
			// 最左侧的六位数字 取0-5位元素，不包括6
			styleNo = sku[0:6]
		} else {
			sku = getCellTrimSpace(f, sheetName, "H", rowindex)
			skuSplit := strings.Split(sku, "-")
			if len(skuSplit) > 1 {
				styleNo = skuSplit[1][0:6]
			}
		}
		if styleNo == "" {
			continue
		}

		// 4. 英文颜色* 取“Color”列下方的内容
		colorEn := getCellTrimSpace(f, sheetName, "F", rowindex)

		// 5. 尺码* 取“Size”列下方的内容
		size := getCellTrimSpace(f, sheetName, "I", rowindex)

		// 6. 订单数量* 取“QTY”列下方的内容
		qtytext := getCellTrimSpace(f, sheetName, "J", rowindex)

		if colorEn == "" || size == "" || qtytext == "" {
			continue
		}

		qtystr := GetDigits(qtytext)
		qty, err = strconv.Atoi(qtystr)
		if err != nil {
			continue
		}

		item := OrderItem{}
		item.PoNo = info.PoNo
		item.StyleNo = styleNo
		item.Qty = qty                                              // 订单数量
		item.Size = size                                            // 尺码
		item.ColorEn = colorEn                                      // 英文颜色
		item.DestCountry = info.DestCountry                         // 目的国
		item.DestPortName = info.DestPortName                       // 目的港
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		info.OrderItems = append(info.OrderItems, item)
	}
	return nil
}
