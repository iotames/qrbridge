package biz

import (
	// "bytes"
	"fmt"
	"strconv"

	// "strings"
	// go get github.com/unidoc/unipdf/v4
	// "github.com/ledongthuc/pdf"
	"github.com/iotames/qrbridge/service"
	"github.com/xuri/excelize/v2"
)

func PoBewcwTransform(inputtpl, inputfile, outputfile string) (info PoInfo, err error) {
	// pdf.DebugOn = true
	// var content string
	// content, err = readPdf(inputfile) // Read local pdf file
	// var ff *os.File
	// ff, _ = os.OpenFile(inputfile, os.O_CREATE|os.O_TRUNC, 0755)
	// defer ff.Close()
	// _, err = io.WriteString(ff, content)
	// fmt.Println(content)

	f, err := service.NewTableFile(inputfile).OpenExcel()
	if err != nil {
		return PoInfo{}, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	sheets := f.GetSheetList()
	// for i, sheet := range sheets {
	poSheetDataParseBewcw(f, sheets[0], &info)
	// }

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

// 从Excel的每个sheet页面解析数据
func poSheetDataParseBewcw(f *excelize.File, sheetName string, info *PoInfo) error {
	var i uint
	for i = 2; i <= 152; i++ {
		item := OrderItem{}
		styleNoStr := getCellTrimSpace(f, sheetName, "A", i) // 提取客户款号
		strno, err := strconv.Atoi(styleNoStr)
		if err != nil {
			continue
		}
		if strno < 1000 {
			continue
		}
		qtytext := getCellTrimSpace(f, sheetName, "I", i)
		qty, err := strconv.Atoi(qtytext)
		if err != nil {
			continue
		}
		item.Qty = qty
		item.StyleNo = fmt.Sprintf("%d", strno) // 客户款号
		item.Desc = getCellTrimSpace(f, sheetName, "B", i)
		item.Size = getCellTrimSpace(f, sheetName, "F", i)
		info.OrderItems = append(info.OrderItems, item)
	}
	return nil
}

// https://github.com/temamagic/rscpdf
// func readPdf(path string) (string, error) {
// 	f, r, err := pdf.Open(path)
// 	// remember close file
// 	if err != nil {
// 		return "", err
// 	}
// 	var buf bytes.Buffer
// 	b, err := r.GetPlainText()
// 	if err != nil {
// 		return "", err
// 	}
// 	buf.ReadFrom(b)
// 	f.Close()
// 	return buf.String(), nil
// }
