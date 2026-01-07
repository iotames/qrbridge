package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoB1ztvTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseB1ztv)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseB1ztv(f *excelize.File, sheetIndex int, info *PoInfo) error {

	var rows [][]string
	var err error
	var rowindex uint
	// var desc, colorEn, size string

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	if sheetIndex == 0 {
		shipSplit := strings.Split(getCellTrimSpace(f, sheetName, "A", 2), "\n")

		if len(shipSplit) > 1 {
			info.DestCountry = strings.TrimSpace(shipSplit[len(shipSplit)-1]) // 目的国
			destPortAddr := shipSplit[len(shipSplit)-2]
			destPortSplit := strings.Split(destPortAddr, " ")
			info.DestPortName = strings.TrimSpace(destPortSplit[len(destPortSplit)-1]) // 目的港

		}
		// fmt.Printf("------shipSplit(%+v)---info.DestCountry(%+v)--info.DestPortName(%+v)--\n", shipSplit, info.DestCountry, info.DestPortName)

		orderSplit := strings.Split(getCellTrimSpace(f, sheetName, "D", 2), "\n")

		for _, order := range orderSplit {
			if strings.Contains(order, "Order Number:") {
				// 客户PO
				info.PoNo = strings.TrimSpace(strings.Replace(order, "Order Number:", "", 1))
			}
			if strings.Contains(order, "Shipment Date:") {
				// 客户交期
				info.DeliveryDateCustomer, _ = time.Parse("02/01/2006", strings.TrimSpace(strings.Replace(order, "Shipment Date:", "", 1)))
			}
		}
		// fmt.Printf("----orderSplit(%+v)--PoNo(%s)--deliveryDateCustomerStr(%s)---\n", orderSplit, info.PoNo, info.DeliveryDateCustomer)
	}
	fmt.Printf("------PoSheetDataParseB1ztv-----info.DestCountry(%+v)--info.DestPortName(%+v)-----\n", info.DestCountry, info.DestPortName)

	if !info.DeliveryDateCustomer.IsZero() {
		// 客户交期
		deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

		// 离厂交期
		deliveryDateFactoryLeave := info.DeliveryDateCustomer.AddDate(0, 0, -7) // 离厂交期=客户交期-7天
		deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")

		// 工厂交期
		deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7) // 工厂交期=离厂交期-7天
		deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
	}

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}
	sizeColIndex := 0
	for i, row := range rows {
		// fmt.Printf("----PoSheetDataParseB1ztv--eachrow(%+v)---\n", row)
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
		if rowindex == 1 {
			// 跳过第一行
			continue
		}

		// 确认尺码列
		if rowindex == 3 || rowindex == 2 {
			for coli, celval := range row {
				if strings.TrimSpace(celval) == "Size" {
					sizeColIndex = coli
				}
			}
		}
		if sizeColIndex < 4 || sizeColIndex > len(row) {
			continue
		}

		// 尺码。G4 OR F3
		size := strings.TrimSpace(row[sizeColIndex])
		if size == "Size" {
			continue
		}

		// 英文颜色。F4 OR E3 列
		colorEn := strings.TrimSpace(row[sizeColIndex-1])

		// 提取客户款号
		// 使用商品标题当客户款号。E4 OR D3 列
		desc := strings.TrimSpace(row[sizeColIndex-2])

		// 订单数量。sheet1 O4 Sheet2 M3 sheet3 O3
		qty := 0
		for ii := sizeColIndex + 1; ii < len(row); ii++ {
			// 尺码列往又数，第一个包含数字的单元格，就是订单数量。
			qtytext := strings.TrimSpace(row[ii])
			if qtytext == "" {
				continue
			}
			qtystr := GetDigits(qtytext)
			qty, err = strconv.Atoi(qtystr)
			if err == nil {
				break
			}
		}
		item := OrderItem{}
		item.PoNo = info.PoNo
		item.StyleNo = desc // 客户款号
		item.Qty = qty      // 订单数量
		item.Desc = desc
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
