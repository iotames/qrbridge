package biz

import (
	"fmt"
	"time"

	"strings"

	"github.com/iotames/qrbridge/service"
	"github.com/xuri/excelize/v2"
)

func PoFileTransform(inputtpl, inputfile, outputfile string) (info PoInfo, err error) {
	f, err := service.NewTableFile(inputfile).OpenExcel()
	if err != nil {
		return PoInfo{}, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	if inputtpl == "Rohnisch" {
		sheets := f.GetSheetList()
		for i, sheet := range sheets {
			poSheetDataParseRohnisch(f, sheet, i, &info)
		}
	}

	err = f.Close()
	if err != nil {
		return PoInfo{}, fmt.Errorf("关闭%s文件失败: %w", inputfile, err)
	}

	// 输出新的EXCEL
	f = service.NewTableFile(outputfile).NewExcel()

	titleRow := []interface{}{"客户款号*", "颜色*", "英文颜色*", "色号", "PO NO*", "尺码*", "工厂交期", "离厂交期*", "客户交期*", "订单数量*", "目的国*"}
	err = f.SetSheetRow("Sheet1", "A1", &titleRow)
	if err != nil {
		return PoInfo{}, fmt.Errorf("fai to SetSheetRow: %w", err)
	}
	for i := 0; i < len(info.OrderItems); i++ {
		rowIndex := i + 2
		err = f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", rowIndex), &[]interface{}{info.OrderItems[i].StyleNo, info.OrderItems[i].Color, info.OrderItems[i].ColorEn, info.OrderItems[i].ColorNo, info.OrderItems[i].PoNo, info.OrderItems[i].Size, info.OrderItems[i].DeliveryDateFactory, info.OrderItems[i].DeliveryDateFactoryLeave, info.OrderItems[i].DeliveryDateCustomer, info.OrderItems[i].Qty, info.OrderItems[i].DestCountry})
		if err != nil {
			return PoInfo{}, fmt.Errorf("fai to SetSheetRow%d: %w", rowIndex, err)
		}
	}
	err = f.SaveAs(outputfile)
	if err != nil {
		return PoInfo{}, fmt.Errorf("保存%s文件失败: %w", outputfile, err)
	}
	err = f.Close()
	return info, err
}

// 从Excel的每个sheet页面解析数据
func poSheetDataParseRohnisch(f *excelize.File, sheetName string, id int, info *PoInfo) error {
	// 第57行为标题行
	// 客户款号C58
	// 款式标题G58 Kay High Waist Tights  Black/Black, XS,  0101
	// noTitle, _ := f.GetCellValue(sheetName, "C57")
	// no, _ := f.GetCellValue(sheetName, "C58")
	poNo := getCellTrimSpace(f, sheetName, "T", 10)                    // 获取PO号
	deliveryDateCustomerTxt := getCellTrimSpace(f, sheetName, "I", 50) // 获取客户交期。I50 25-09-15
	deliveryDateCustomerStr := "20" + deliveryDateCustomerTxt          // 客户交期
	deliveryDateFactoryLeaveStr := ""                                  // 离厂交期
	deliveryDateFactoryStr := ""                                       // 工厂交期
	deliveryDateCustomer, err := time.Parse("2006-01-02", deliveryDateCustomerStr)
	if err == nil {
		deliveryDateFactoryLeave := deliveryDateCustomer.AddDate(0, 0, -7)
		deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, 7)
		deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")
		deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
	}
	destCountry := getCellTrimSpace(f, sheetName, "F", 31) // 目的国
	ids := getOkOrderItemRowIndexs(f, sheetName, "C", 58, 5, []string{"No."})
	for _, rowindex := range ids {
		item := OrderItem{}
		item.StyleNo = getCellTrimSpace(f, sheetName, "C", rowindex) // 提取客户款号
		item.Desc = getCellTrimSpace(f, sheetName, "G", rowindex)
		item.ColorEn, item.Size = getRohnischColorSizeByDesc(item.Desc) // 获取颜色尺码
		item.PoNo = poNo
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。客户交期-7天
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		qtyStr := getCellTrimSpace(f, sheetName, "AA", rowindex)    // 原始数据。订单数量。必填
		fmt.Sscanf(qtyStr, "%d", &item.Qty)                         // 转换为整型。订单数量。必填
		item.DestCountry = destCountry
		info.OrderItems = append(info.OrderItems, item)
		fmt.Printf("----sheet(%d-%s)---rowindex(%d)---orderItem(%+v)------\n", id, sheetName, rowindex, item)
	}
	return nil
}

func getCellTrimSpace(f *excelize.File, sheetName string, col string, rowIndex uint) string {
	cell := fmt.Sprintf("%s%d", col, rowIndex)
	cellValue, _ := f.GetCellValue(sheetName, cell)
	cellValue = strings.TrimSpace(cellValue)
	return cellValue
}

// 获取有效的订单条目行索引
func getOkOrderItemRowIndexs(f *excelize.File, sheetName string, col string, startRowIndex uint, setplen uint8, exclude []string) []uint {
	cell := fmt.Sprintf("%s%d", col, startRowIndex)
	cellValue, _ := f.GetCellValue(sheetName, cell)
	cellValue = strings.TrimSpace(cellValue)
	var okindexs []uint
	var addIndex uint
	var trycount uint8
	for {
		skip := false
		if trycount > setplen {
			// 找了好多行，确实没有效数据了。结束
			break
		}

		if cellValue == "" {
			// 单元格等于空值，不符合要求，跳过继续寻找下一行
			// fmt.Println("---Skip-getOkOrderItemRowIndex--emptyValue--", sheetName, col, startRowIndex+addIndex)
			nextCellValue(f, sheetName, col, startRowIndex, &cellValue, &addIndex)
			trycount += 1
			continue
		}

		for _, v := range exclude {
			if cellValue == strings.TrimSpace(v) {
				// 单元格等于某个值，不符合要求，跳过继续寻找下一行
				fmt.Printf("---Skip-getOkOrderItemRowIndex--Sheet(%s)-cel(%s%d)(%s)--\n", sheetName, col, startRowIndex+addIndex, cellValue)
				skip = true
				break
			}
		}
		if skip {
			nextCellValue(f, sheetName, col, startRowIndex, &cellValue, &addIndex)
			trycount += 1
			continue
		}

		// 找到有效值，加入结果集
		okindexs = append(okindexs, startRowIndex+addIndex)
		nextCellValue(f, sheetName, col, startRowIndex, &cellValue, &addIndex)
		trycount = 0
	}
	return okindexs
}

func nextCellValue(f *excelize.File, sheetName string, col string, startRowIndex uint, cellValue *string, addIndex *uint) {
	*addIndex += 1
	cell := fmt.Sprintf("%s%d", col, startRowIndex+*addIndex)
	*cellValue, _ = f.GetCellValue(sheetName, cell)
	*cellValue = strings.TrimSpace(*cellValue)
}

func getRohnischColorSizeByDesc(desc string) (color, size string) {
	descsplit := strings.Split(desc, ",")
	splitlen := len(descsplit)
	if splitlen > 2 {
		size = strings.TrimSpace(descsplit[splitlen-2]) // 倒数第二个
		titleColor := strings.TrimSpace(descsplit[0])
		colorsplit := strings.Split(titleColor, " ")
		colorsplitlen := len(colorsplit)
		if colorsplitlen > 1 {
			color = strings.TrimSpace(colorsplit[colorsplitlen-1])
		}
	}
	return
}
