package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoAh8swTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, 0, PoSheetDataParseAh8sw)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseAh8sw(f *excelize.File, sheetIndex int, info *PoInfo) error {

	var rows [][]string
	var err error
	var rowindex uint
	var qty int

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)
	// 备注： PO 未体现色号，工厂交期，目的国，目的港

	if sheetIndex == 0 {
		// 2025/10/30  K7
		deliveryDateCustomerText := getCellTrimSpace(f, sheetName, "K", 7)
		// 客户交期
		info.DeliveryDateCustomer, err = time.Parse("2006/01/02", deliveryDateCustomerText)
		if err != nil {
			// "02-Jan-06"对应：日-月缩写-年（2位）
			info.DeliveryDateCustomer, _ = time.Parse("02-Jan-06", deliveryDateCustomerText)
		}

		info.PoNo = getCellTrimSpace(f, sheetName, "B", 5) // 客户PO B5 == M10
		fmt.Printf("----PoNo(%s)----deliveryDateCustomerText(%s)--deliveryDateCustomerStr(%s)---\n", info.PoNo, deliveryDateCustomerText, info.DeliveryDateCustomer)
	}
	fmt.Printf("------PoSheetDataParseAh8sw-----info.DestCountry(%+v)--info.DestPortName(%+v)-----\n", info.DestCountry, info.DestPortName)

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
		if rowindex < 10 {
			// 跳过前9行
			continue
		}

		// 客户款号。B10 开始
		styleNo := getCellTrimSpace(f, sheetName, "B", rowindex)
		// 英文颜色。E10开始
		colorEn := getCellTrimSpace(f, sheetName, "E", rowindex)
		// 色号。C10开始
		colorNo := getCellTrimSpace(f, sheetName, "C", rowindex)
		// 尺码。H10开始
		size := getCellTrimSpace(f, sheetName, "H", rowindex)
		// 订单数量。J10开始
		qtytext := getCellTrimSpace(f, sheetName, "J", rowindex)
		if styleNo == "" || colorEn == "" || size == "" || qtytext == "" {
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
		item.ColorNo = colorNo
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
