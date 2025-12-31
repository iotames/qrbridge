package biz

import (
	"fmt"
	"strings"

	"github.com/iotames/qrbridge/service"
	"github.com/xuri/excelize/v2"
)

func potransform(inputfile, outputfile string, sheetIndex int, transfunc func(f *excelize.File, sheetIndex int, info *PoInfo) error) (info PoInfo, err error) {
	f, err := service.NewTableFile(inputfile).OpenExcel()
	if err != nil {
		return PoInfo{}, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	sheets := f.GetSheetList()
	sheetLen := len(sheets)
	if sheetLen == 0 {
		return PoInfo{}, fmt.Errorf("没有sheet页")
	}
	if sheetIndex > sheetLen-1 {
		return PoInfo{}, fmt.Errorf("sheet页索引超限(now%d/max%d)", sheetIndex, sheetLen-1)
	}
	if sheetIndex < 0 {
		for si, _ := range sheets {
			transfunc(f, si, &info)
		}
	} else {
		transfunc(f, sheetIndex, &info)
	}
	err = f.Close()
	if err != nil {
		return PoInfo{}, fmt.Errorf("关闭%s文件失败: %w", inputfile, err)
	}
	err = poOutputExcel(outputfile, info)
	if err != nil {
		return PoInfo{}, fmt.Errorf("输出Excel文件失败: %w", err)
	}
	return info, err
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
