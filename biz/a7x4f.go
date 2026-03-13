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
	var dateindex uint // 客户交期列
	// var qtyindex uint // Qty per size列

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期

	info.DestCountry = "Poland"
	info.DestPortName = "Gdansk"

	sheetName := f.GetSheetName(sheetIndex)

	// 获取PO号
	poNoText := getCellTrimSpace(f, sheetName, "C", 2)
	if strings.Contains(poNoText, " ") {
		poNoSplit := strings.Split(poNoText, " ")
		if len(poNoSplit) > 1 {
			info.PoNo = poNoSplit[len(poNoSplit)-1]
		}
	}

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}
	fmt.Printf("---sheet(%s)---获取总行数(%d)--\n", sheetName, len(rows))
	nowStyleNo := ""
	for i, row := range rows {
		// fmt.Printf("----PoSheetDataParseAh8sw--eachrow(%+v)---\n", row)
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// for _, cell := range row {
		// 	if strings.TrimSpace(cell) == "2025-09-23" {
		// 		fmt.Printf("-----sheetName(%s)-rowindex(%d)--row(%s)---\n", sheetName, rowindex, strings.Join(row, "|"))
		// 	}
		// }

		// // 跳出空数据行
		// if strings.TrimSpace(row[0]) == "" {
		// 	continue
		// }

		// 客户款号
		styleNo := strings.TrimSpace(row[0])

		deliveryDateCustomerText := getCellTrimSpace(f, sheetName, "D", rowindex)

		// 1. 客户交期 2026/4/30, 04-30-26
		info.DeliveryDateCustomer, err = time.Parse("2006-01-02", deliveryDateCustomerText)
		// if err != nil {
		// 	// "02-Jan-06"对应：日-月缩写-年（2位）
		// 	// "01-02-06" 对应：月-日-年（2位）
		// 	info.DeliveryDateCustomer, err = time.Parse("2006-01-02", deliveryDateCustomerText)
		// }
		if err != nil {
			// 取不到客户交期，为无效列，跳过。
			continue
		}
		// if !info.DeliveryDateCustomer.IsZero() {}
		// 客户交期
		deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

		// 离厂交期
		deliveryDateFactoryLeave := info.DeliveryDateCustomer.AddDate(0, 0, -7) // 离厂交期=客户交期-7天
		deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")

		// 工厂交期
		deliveryDateFactory := deliveryDateFactoryLeave // deliveryDateFactoryLeave.AddDate(0, 0, -7) // 工厂交期=离厂交期-7天
		deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")

		if styleNo != "" {
			nowStyleNo = styleNo
		}

		qtytext := ""
		// 获取客户交期的列索引dateindex
		for ci, cv := range row {
			if strings.TrimSpace(cv) == deliveryDateCustomerText {
				dateindex = uint(ci)
				break
			}
		}
		// Qty persize 计算具体尺码的数量索引
		// qtytext = row[dateindex+6]

		fmt.Printf("-----sheetName(%s)-rowindex(%d)--dateText(%s)-qtytext(%s)--\n", sheetName, rowindex, deliveryDateCustomerText, qtytext)

		// // 获取订单数量
		// qtystr := GetDigits(qtytext)
		// qty, err = strconv.Atoi(qtystr)
		// if err != nil {
		// 	continue
		// }

		size := strings.TrimSpace(row[dateindex+2])
		colorEnText := strings.TrimSpace(row[dateindex-1])
		colorSplit := strings.Split(colorEnText, " ")
		colorEn := ""
		colorNo := ""
		if len(colorSplit) == 1 {
			colorSplit = strings.Split(colorEnText, "\n")
		}

		if len(colorSplit) > 1 {
			colorEn = strings.Replace(strings.Join(colorSplit[1:], " "), "\n", "", -1)
			colorNo = strings.TrimSpace(colorSplit[0])
		}
		if colorEn == "" || colorNo == "" {
			colorEn = colorEnText
			colorNo = colorEnText
		}

		fmt.Printf("--sheetName(%s)-rowindex(%d)--row(%s)-nowStyleNo(%s)-qty(%d)--colorEn(%s-%s-%d])-size(%s)--poNo(%s)-deliveryDateCustomerText(%s)--deliveryDateCustomerStr(%s)---\n",
			sheetName, i, strings.Join(row, "|"), nowStyleNo, qty, colorEn, colorEnText, len(colorSplit), size, info.PoNo, deliveryDateCustomerText, info.DeliveryDateCustomer)

		for count := 0; count < 10; count++ {
			childRow := rows[i+count]
			if len(childRow) < 8 {
				break
			}
			sizeVal := strings.TrimSpace(childRow[dateindex+2])
			if sizeVal == "0" {
				sizeVal = strings.TrimSpace(childRow[dateindex+3])
			}
			if sizeVal == "" {
				break
			}
			// Qty persize 计算具体尺码的数量索引
			qtytext = childRow[dateindex+6]
			if qtytext == "" {
				qtytext = childRow[dateindex+7]
			}
			// 获取订单数量
			qtystr := GetDigits(qtytext)
			qty, err = strconv.Atoi(qtystr)
			if err != nil {
				break
			}

			item := OrderItem{}
			item.PoNo = info.PoNo
			item.StyleNo = nowStyleNo
			item.Qty = qty      // 订单数量
			item.Size = sizeVal // 尺码
			item.ColorNo = colorNo
			item.ColorEn = colorEn                                      // 英文颜色
			item.DestCountry = info.DestCountry                         // 目的国
			item.DestPortName = info.DestPortName                       // 目的港
			item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
			item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
			item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
			info.OrderItems = append(info.OrderItems, item)
			fmt.Printf("---child-rowindex(%d)--nowStyleNo(%s)-qty(%d)--colorEn(%s)--sizeVal(%s)--\n", i+count, nowStyleNo, qty, colorEn, sizeVal)

		}

	}
	return nil
}
