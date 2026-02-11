package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA86sgTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseA86sg)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseA86sg(f *excelize.File, sheetIndex int, info *PoInfo) error {
	var rows [][]string
	var err error
	var rowindex uint
	// var qty int
	// var desc, colorEn, size string

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}

	// 目的港
	info.DestPortName = "VALENCIA"

	// startOrderData := false

	// var sizeList []string

	var styleNo string
	var deliveryDateCustomerTxt string
	var colorMap map[string]map[string]int
	var lastColorEn string
	var lastColorRowIndex uint
	for i, row := range rows {
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		rowDone := false
		var okrow []string
		for _, cell := range row {
			cell := strings.TrimSpace(cell)
			if cell == "" {
				continue
			}
			okrow = append(okrow, cell)
			if rowindex == 2 && strings.Contains(cell, "Style No:") {
				// 1、C2内容包含“Style No: ”的sheet为详情页
				// 2、“Style No: ”右侧的内容为客户款号
				styleNo = strings.TrimSpace(strings.Replace(cell, "Style No:", "", 1))
				rowDone = true
				break
			}
			// 取每个详情页“ETD:  ”右侧的日期
			if rowindex == 6 && strings.Contains(cell, "ETD:") {
				deliveryDateCustomerTxt = strings.TrimSpace(strings.Replace(cell, "ETD:", "", 1))
				info.DeliveryDateCustomer, _ = time.Parse("2006-01-02", deliveryDateCustomerTxt)
			}
			// 英文颜色
			if rowindex > 8 && strings.HasPrefix(cell, "Color:") {
				colorText := strings.TrimSpace(strings.Replace(cell, "Color:", "", 1))
				colorSplit := strings.Split(colorText, " ")
				if len(colorSplit) > 1 {
					colorEn := strings.Join(colorSplit[1:], " ")
					if colorMap == nil {
						colorMap = make(map[string]map[string]int)
					}
					colorMap[colorEn] = make(map[string]int)
					lastColorEn = colorEn
					lastColorRowIndex = rowindex
				}
			}
			// 颜色对应的尺码
			if rowindex == lastColorRowIndex+1 {
				if cell == "Total" {
					break
				}
				colorMap[lastColorEn][cell] = 0
			}
			// 颜色的具体尺码的数量 TODO
			if rowindex == lastColorRowIndex+3 {
				if cell == "Total" {
					continue
				}
				qtyStr := GetDigits(cell)
				qty, _ := strconv.Atoi(qtyStr)

				colorMap[lastColorEn]["xxxx"] = qty
			}

		}
		if rowDone {
			continue
		}

		fmt.Println("---okrow(%+v)----\n", strings.Join(okrow, "|"))

		if !info.DeliveryDateCustomer.IsZero() {
			// 客户交期
			deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

			// 离厂交期
			deliveryDateFactoryLeave := info.DeliveryDateCustomer.AddDate(0, 0, -8) // 离厂交期=客户交期-8天
			deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")

			// 工厂交期
			deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7) // 工厂交期=离厂交期-7天
			deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
		}
		fmt.Printf("------PoSheetDataParseA86sg----DeliveryDateCustomer(%+v)--info.DestPortName(%+v)-----\n", info.DeliveryDateCustomer, info.DestPortName)

		// 跳出空数据行
		vala := strings.TrimSpace(row[0])
		if vala == "" {
			continue
		}

		// if startOrderData {
		// 				qtystr = strings.TrimSpace(qtystr)
		// 			size := sizeList[iii]
		// 			qtystr = GetDigits(qtystr)
		// 			qty, _ = strconv.Atoi(qtystr)

		// 			item := OrderItem{}
		// 			item.PoNo = info.PoNo
		// 			item.StyleNo = styleNo                                      // 客户款号
		// 			item.ColorNo = colorNo                                      // 色号
		// 			item.Qty = qty                                              // 订单数量
		// 			item.Size = size                                            // 尺码
		// 			item.ColorEn = colorEn                                      // 英文颜色
		// 			item.DestCountry = info.DestCountry                         // 目的国
		// 			item.DestPortName = info.DestPortName                       // 目的港
		// 			item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		// 			item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
		// 			item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		// 			info.OrderItems = append(info.OrderItems, item)

		// }

	}
	return nil
}
