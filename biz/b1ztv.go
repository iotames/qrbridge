package biz

import (
	"github.com/xuri/excelize/v2"
)

func PoB1ztvTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseB1ztv)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseB1ztv(f *excelize.File, sheetIndex int, info *PoInfo) error {
	return nil
}
