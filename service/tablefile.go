package service

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	// _ "golang.org/x/image/webp"
	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"
)

type TableFile struct {
	filePath  string
	excelFile *excelize.File
}

func NewTableFile(filepath string) *TableFile {
	return &TableFile{filePath: filepath}
}

func (f *TableFile) OpenExcel() (ef *excelize.File, err error) {
	f.excelFile, err = excelize.OpenFile(f.filePath)
	return f.excelFile, err
}

func (f *TableFile) NewExcel() (ef *excelize.File) {
	f.excelFile = excelize.NewFile()
	return f.excelFile
}

// SetRowData. rowi startbegin 1
func SetRowDataByExcel(filepath, sheetName string, data []interface{}, rowi int) error {
	f, err := excelize.OpenFile(filepath)
	if rowi < 1 {
		panic("Error in ExcelService.SetRowData: arg rowi must greater than 0")
	}
	coli := 'A'
	for _, cellValue := range data {
		axis := fmt.Sprintf("%c%d", coli, rowi)
		if coli > 'Z' {
			add := coli - 'Z'
			coll := 'A' + add - 1
			axis = fmt.Sprintf("A%c%d", coll, rowi)
		}
		f.SetCellValue(sheetName, axis, cellValue)
		coli++
	}
	return err
}
