package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA7x4fTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseA7x4f)
}

// 从Excel的多个sheet页面解析数据
func PoSheetDataParseA7x4f(f *excelize.File, sheetIndex int, info *PoInfo) error {

	var rows [][]string
	var err error
	var rowindex uint
	var qty int

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期

	// info.DestCountry = ""
	// info.DestPortName = ""

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

		info.PoNo = getCellTrimSpace(f, sheetName, "A", rowindex)

		styleNo := ""

		if styleNo == "" {
			continue
		}

		qtytext := ""

		qtystr := GetDigits(qtytext)
		qty, err = strconv.Atoi(qtystr)
		if err != nil {
			continue
		}

		item := OrderItem{}
		item.PoNo = info.PoNo
		item.StyleNo = styleNo
		item.Qty = qty                                              // 订单数量
		item.Size = ""                                              // 尺码
		item.ColorEn = ""                                           // 英文颜色
		item.DestCountry = info.DestCountry                         // 目的国
		item.DestPortName = info.DestPortName                       // 目的港
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		info.OrderItems = append(info.OrderItems, item)
	}
	return nil
}
