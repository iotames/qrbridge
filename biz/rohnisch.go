package biz

import (
	"fmt"

	"strings"

	"github.com/xuri/excelize/v2"
)

func PoSheetRohnisch(f *excelize.File, sheetName string, id int) PoInfo {
	// 第57行为标题行
	// 客户款号C58
	// 款式标题G58 Kay High Waist Tights  Black/Black, XS,  0101
	// noTitle, _ := f.GetCellValue(sheetName, "C57")
	// no, _ := f.GetCellValue(sheetName, "C58")
	ids := getOkOrderItemRowIndexs(f, sheetName, "C", 58, 5, []string{"No."})
	info := PoInfo{OrderItems: []OrderItem{}}
	for _, rowindex := range ids {
		item := OrderItem{}
		item.StyleNo = getCellTrimSpace(f, sheetName, "C", rowindex)
		item.Desc = getCellTrimSpace(f, sheetName, "G", rowindex)
		item.Color, item.Size = getRohnischColorSizeByDesc(item.Desc)
		info.OrderItems = append(info.OrderItems, item)
		fmt.Printf("----sheet(%d-%s)---rowindex(%d)---orderItem(%+v)------\n", id, sheetName, rowindex, item)
	}
	return info
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
			fmt.Println("---Skip-getOkOrderItemRowIndex--emptyValue--", sheetName, col, startRowIndex+addIndex)
			nextCellValue(f, sheetName, col, startRowIndex, &cellValue, &addIndex)
			trycount += 1
			continue
		}

		for _, v := range exclude {
			if cellValue == strings.TrimSpace(v) {
				// 单元格等于某个值，不符合要求，跳过继续寻找下一行
				fmt.Println("---Skip-getOkOrderItemRowIndex--", sheetName, col, startRowIndex+addIndex)
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
