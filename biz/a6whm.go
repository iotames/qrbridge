package biz

import (
	// "archive/zip"
	// "encoding/xml"
	// "bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA6whmTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseA6whm)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseA6whm(f *excelize.File, sheetIndex int, info *PoInfo) error {

	// 客户	    业务员	客户款号*	颜色*	英文颜色*	色号	PO NO*	尺码*	工厂交期	离厂交期*	客户交期*	订单数量*	目的国*	目的港	其他
	// A6WHM	张真真	有	无	有	无	有	有	离厂交期-7天	等于客户交期	有	有	有	无
	// 1、一个Excel对应一个PO，然后一个合同由多个PO组成

	var rows [][]string
	var err error
	// var poInt int
	// var rowindex uint

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	if sheetIndex == 0 {
		// 目的国
		info.DestCountry = "Netherland" // getCellTrimSpace(f, sheetName, "E", 15) // THE NETHERLANDS
		info.DestPortName = "阿姆斯特丹"
		ponoStr := getCellTrimSpace(f, sheetName, "J", 11)
		_, err := strconv.Atoi(ponoStr)
		// PO NO 无法转换成整数，则不保存。
		if err == nil {
			info.PoNo = ponoStr
		}
		// 客户交期 日.月.年
		deliveryDateCustomerTxt := getCellTrimSpace(f, sheetName, "M", 11)

		info.DeliveryDateCustomer, _ = time.Parse("02.01.2006", deliveryDateCustomerTxt)
	}

	if !info.DeliveryDateCustomer.IsZero() {
		// 客户交期
		deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

		// 离厂交期
		deliveryDateFactoryLeave := info.DeliveryDateCustomer // 离厂交期==客户交期
		deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")

		// 工厂交期
		deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7) // 工厂交期=离厂交期-7天
		deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
	}

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}
	for _, row := range rows {
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		// rowindex = uint(i + 1)
		// 获取当前行的第一列数据
		rowA := strings.TrimSpace(row[0])

		// 提取客户款号
		// isOrderDetailRow
		orderId, err := strconv.Atoi(rowA)
		if err != nil {
			continue
		}
		if orderId < 1000000 {
			continue
		}
		// for coli, celval := range row {
		// 	fmt.Printf("---parse-each-row-col--coli(%d)-----celval(%s)-------\n", coli, celval)
		// }
		// fmt.Printf("-----PoSheetDataParseA6whm----parse-each-row--row(%+v)---\n", row)
		tidyrow := []string{} // len = 6
		for _, celval := range row {
			tidycell := strings.TrimSpace(celval)
			if tidycell != "" {
				tidyrow = append(tidyrow, tidycell)
			}
		}
		if len(tidyrow) < 4 {
			continue
		}
		// 取出完整的商品标题或者描述

		// 情况1：商品描述只占1个单元格。直接取出。
		desc := tidyrow[1]
		qtytext := tidyrow[2]
		qtystr := GetDigits(qtytext)
		qty, err := strconv.Atoi(qtystr)
		if err != nil {
			// 情况2：商品描述占2个单元格
			qtytext = tidyrow[3]
			qtystr = GetDigits(qtytext)
			qty, err = strconv.Atoi(qtystr)
			desc1 := tidyrow[1]
			// A. 先对第1个单元格的字符串进行换行符分隔。
			descSplit := strings.Split(desc1, "\n")
			addStr := ""

			if len(descSplit) > 1 {
				addStr = descSplit[1]
			}
			desc2 := tidyrow[2]
			// B. 如果分隔成功，取出两部分。分隔的第1部分和第2个单元格的商品描述拼接。分隔的第2部分拼接到末尾。
			// B. 如果分隔失败，则效果等同把两个单元格直接拼接。
			desc = descSplit[0] + desc2 + addStr
		}
		// 在完整的商品标题或者描述中，通过逗号(,)分隔符，取出颜色和尺码。
		colorEn, size := getA6whmColorSizeByDesc(desc)

		// fmt.Printf("-----PoSheetDataParseA6whm--each-row--tidyrow(%+v)-qtytext(%s)-qtystr(%s)---\n", tidyrow, qtytext, qtystr)

		if err != nil {
			continue
		}
		item := OrderItem{}
		item.PoNo = info.PoNo
		item.StyleNo = fmt.Sprintf("%d", orderId) // 客户款号
		item.Qty = qty                            // 订单数量
		item.Desc = desc
		item.Size = size       // 尺码
		item.ColorEn = colorEn // 英文颜色

		item.DestCountry = info.DestCountry // 目的国
		item.DestPortName = info.DestPortName
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		info.OrderItems = append(info.OrderItems, item)
	}
	return nil
}

func getA6whmColorSizeByDesc(desc string) (color, size string) {
	descsplit := strings.Split(desc, ",")
	if len(descsplit) < 3 {
		return
	}
	return strings.TrimSpace(descsplit[2]), strings.TrimSpace(descsplit[1])
}
