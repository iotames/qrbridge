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

type SizeQty struct {
	Index            int
	ColorNo, ColorEn string
	Size             string
	Qty              int
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseA86sg(f *excelize.File, sheetIndex int, info *PoInfo) error {
	var rows [][]string
	var err error
	var rowindex uint
	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}
	// fmt.Printf("----sheet(%s)--rows.len(%d)----\n", sheetName, len(rows))

	// 目的港
	info.DestPortName = "Aarhus"
	// 目的国
	info.DestCountry = "Denmark"

	var styleNo string
	var deliveryDateCustomerTxt string
	var sizeQtyListTemp []SizeQty
	var totalSizeQtyList []SizeQty
	var nowColorEn, nowColorNo string
	var currentColorRowIndex uint
	for i, row := range rows {
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		rowDone := false
		// var okrow []string
		sizeIndex := 0
		// fmt.Printf("----sheet(%s)--rowindex(%d)----row(%s)-----\n", sheetName, rowindex, strings.Join(row, "|"))
		for _, cell := range row {
			cell := strings.TrimSpace(cell)
			if cell == "" {
				continue
			}
			// ---row(Purchase Order no: APO2128)---
			if rowindex == 1 {
				if strings.Contains(cell, "Purchase Order no:") {
					info.PoNo = strings.TrimSpace(strings.Replace(cell, "Purchase Order no:", "", 1))
					rowDone = true
					break
				}
			}
			// okrow = append(okrow, cell)
			if rowindex == 2 {
				if strings.Contains(cell, "Style No:") {
					// 1、C2内容包含“Style No: ”的sheet为详情页
					// 2、“Style No: ”右侧的内容为客户款号
					styleNo = strings.TrimSpace(strings.Replace(cell, "Style No:", "", 1))
					rowDone = true
					break
				} else {
					break
				}
			}
			// 取每个详情页“ETD:  ”右侧的日期
			if rowindex == 6 && strings.Contains(cell, "ETD:") {
				deliveryDateCustomerTxt = strings.TrimSpace(strings.Replace(cell, "ETD:", "", 1))
				info.DeliveryDateCustomer, _ = time.Parse("2006-01-02", deliveryDateCustomerTxt)
			}
			// 英文颜色
			// 有时候在第9行，有时候又7第
			if rowindex > 6 && strings.HasPrefix(cell, "Color:") {
				// 去掉Color:
				colorText := strings.TrimSpace(strings.Replace(cell, "Color:", "", 1))
				// 空格分隔成数组
				colorSplit := strings.Split(colorText, " ")

				if len(colorSplit) > 1 {
					// 取数组第一部分为色号
					nowColorNo = colorSplit[0]
					// 取数组第一部分后的内容，还原成字符串，为英文颜色。
					nowColorEn = strings.Join(colorSplit[1:], " ")
					currentColorRowIndex = rowindex
				}
				// fmt.Printf("----foreach-cell-rowindex(%d)--cell(%s)--colorSplit(%+v)--nowColorEn(%s)--currentColorRowIndex(%d)-\n", rowindex, cell, colorSplit, nowColorEn, currentColorRowIndex)
			}
			// 颜色对应的尺码
			if rowindex == currentColorRowIndex+1 && nowColorEn != "" {
				if cell == "Total" {
					break
				}
				// 初始化尺码数据
				sizeQtyListTemp = append(sizeQtyListTemp, SizeQty{Index: sizeIndex, ColorEn: nowColorEn, ColorNo: nowColorNo, Size: cell})
				sizeIndex++
			}

			// 不同颜色不同尺码的数量
			if rowindex == currentColorRowIndex+3 && len(sizeQtyListTemp) > 0 {
				if cell == "Total" {
					continue
				}
				qtyStr := GetDigits(cell)
				qty, _ := strconv.Atoi(qtyStr)
				sizeQtyListTemp[sizeIndex].Qty = qty
				sizeIndex++
				// fmt.Printf("----foreach--row--cell(%s)--sizeIndex(%d)--sizeQtyListTemp.len(%d)--\n", cell, sizeIndex, len(sizeQtyListTemp))
				if sizeIndex == len(sizeQtyListTemp) {
					totalSizeQtyList = append(totalSizeQtyList, sizeQtyListTemp...)
					// 清空尺码数量
					sizeQtyListTemp = []SizeQty{}
					break
				}
			}
			// fmt.Printf("---------sizeQtyListTemp(%+v)-------\n", sizeQtyListTemp)
		}
		if rowDone {
			continue
		}

		// fmt.Printf("---rowindex(%d)--okrow(%s)----\n", rowindex, strings.Join(okrow, "|"))

		// // 跳出空数据行
		// vala := strings.TrimSpace(row[0])
		// if vala == "" {
		// 	continue
		// }

	}

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

	for _, sizeQty := range totalSizeQtyList {
		item := OrderItem{}
		item.PoNo = info.PoNo
		item.StyleNo = styleNo                                      // 客户款号
		item.ColorNo = sizeQty.ColorNo                              // 色号
		item.Qty = sizeQty.Qty                                      // 订单数量
		item.Size = sizeQty.Size                                    // 尺码
		item.ColorEn = sizeQty.ColorEn                              // 英文颜色
		item.DestCountry = info.DestCountry                         // 目的国
		item.DestPortName = info.DestPortName                       // 目的港
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		info.OrderItems = append(info.OrderItems, item)
	}

	return nil
}
