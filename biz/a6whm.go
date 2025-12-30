package biz

import (
	// "archive/zip"
	// "encoding/xml"
	// "bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func PoA6wHmTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, 0, PoSheetDataParseA6wHm)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseA6wHm(f *excelize.File, sheetName string, info *PoInfo) error {
	// 离厂交期==客户交期
	// 工厂交期==离厂交期-7天
	// 客户	业务员	客户款号*	颜色*	英文颜色*	色号	PO NO*	尺码*	工厂交期	离厂交期*	客户交期*	订单数量*	目的国*	目的港	其他
	// A6WHM	张真真	有	无	有	无	有	有	离厂交期-7天	等于客户交期	有	有	有	无	1、一个Excel对应一个PO，然后一个合同由多个PO组成

	var rows [][]string
	var err error
	// var poInt int
	var rowindex uint
	var destCountry string // 目的国
	// var deliveryDateCustomer time.Time
	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期

	destCountryText := getCellTrimSpace(f, sheetName, "B", 14)
	destCountrySplit := strings.Split(destCountryText, "&")
	if len(destCountrySplit) > 0 {
		destCountry = strings.TrimSpace(destCountrySplit[0])
	}

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}
	for i, row := range rows {
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// 获取当前行的第一列数据rowB
		rowB := strings.TrimSpace(row[1])
		// fmt.Printf("----rowindex(%d)--rowB(%s)--\n", rowindex, rowB)

		// // 当前rowB数据是PO编号，设置正整数的PO号，然后继续解析下一行
		// if strings.HasPrefix(rowB, "PO") {
		// 		postr := strings.TrimSpace(strings.Replace(strings.Replace(rowB, "No.", "", 1), "PO", "", 1))
		// 		// 设置PO号
		// 		poInt, err = strconv.Atoi(postr)
		// 		if err == nil {
		// 			// PO号提取成功
		// 			continue
		// 		}
		// }

		// // 当前rowA数据包含Shipment Date
		// if strings.Contains(rowB, "Shipment Date") {
		// 	if len(row) > 2 {
		// 		dateText := row[2]
		// 		datesplit := strings.Split(dateText, "\n")
		// 		if len(datesplit) > 1 {
		// 			deliveryDateCustomerStr = "20" + datesplit[1] // 客户交期

		// 			deliveryDateCustomer, err = time.Parse("2006-01-02", deliveryDateCustomerStr)
		// 			if err == nil {
		// 				// 客户交期的前7天是离厂交期，离厂交期的前5天是工厂交期。
		// 				deliveryDateFactoryLeave := deliveryDateCustomer.AddDate(0, 0, -7) // 离厂交期。必填。客户交期-7天
		// 				deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7)  // 工厂交期。非必填。离厂交期-7天
		// 				deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")
		// 				deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
		// 			}
		// 		}
		// 	}
		// }

		// // 跳过订单详情区块的标题行
		// if rowB == "No." {
		// 	continue
		// }

		// 提取客户款号
		// isOrderDetailRow
		orderId, err := strconv.Atoi(rowB)
		if err != nil {
			continue
		}
		if orderId < 30000000 {
			continue
		}
		fmt.Printf("-----PoSheetDataParseA6wHm--row(%+v)---\n", row)
		// 商品标题或者详情
		desc := getCellTrimSpace(f, sheetName, "D", rowindex)
		colorEn, size := getA6wHmColorSizeByDesc(desc)

		// 提取订单数量
		qtytext := getCellTrimSpace(f, sheetName, "G", rowindex)
		qtytext = strings.Replace(qtytext, "Piece", "", 1)
		qtystr := strings.TrimSpace(strings.Replace(qtytext, ",", "", 1))
		qty, err := strconv.Atoi(qtystr)
		if err != nil {
			continue
		}

		// for coli, celval := range row {}
		// TODO "目的港"
		item := OrderItem{}
		// item.PoNo = fmt.Sprintf("PO%d", poInt)
		item.StyleNo = fmt.Sprintf("%d", orderId) // 客户款号
		item.Qty = qty                            // 订单数量
		item.Desc = desc
		item.Size = size       // 尺码
		item.ColorEn = colorEn // 英文颜色

		item.DestCountry = destCountry                              // 目的国
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。客户交期-7天
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-5天
		info.OrderItems = append(info.OrderItems, item)
	}
	return nil
}

func getA6wHmColorSizeByDesc(desc string) (color, size string) {
	descsplit := strings.Split(desc, ",")
	if len(descsplit) < 3 {
		return
	}
	return strings.TrimSpace(descsplit[2]), strings.TrimSpace(descsplit[1])
}
