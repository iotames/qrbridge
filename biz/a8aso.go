package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/xuri/excelize/v2"
)

func PoA8asoTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseA8aso)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseA8aso(f *excelize.File, sheetIndex int, info *PoInfo) error {
	var rows [][]string
	var err error
	// var rowindex uint
	var qty int
	// var desc, colorEn, size string

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	// PO 固定取H3
	potxt := getCellTrimSpace(f, sheetName, "H", 3)
	potxt = strings.Replace(potxt, `No.`, ``, 1)
	potxt = strings.TrimSpace(potxt)
	if potxt != "" {
		info.PoNo = potxt
	}
	// 目的国*	无

	// 客户交期*	固定取C6
	deliveryDateCustomerTxt := getCellTrimSpace(f, sheetName, "C", 6)
	deliveryDateCustomerTxtSplit := strings.Split(deliveryDateCustomerTxt, "\n")
	if len(deliveryDateCustomerTxtSplit) > 0 {
		deliveryDateCustomerTxt = strings.TrimSpace(deliveryDateCustomerTxtSplit[0])
	}
	info.DeliveryDateCustomer, _ = time.Parse("2006-01-02", deliveryDateCustomerTxt)

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

	// 目的港 取首个sheet页A列最后一行，“,”左边的字母，不要数字
	listaval := []string{}
	for _, tmprv := range rows {
		if len(tmprv) == 0 {
			continue
		}
		tmprva := strings.TrimSpace(tmprv[0])
		if tmprva != "" && strings.Contains(tmprva, ",") {
			listaval = append(listaval, tmprva)
		}
	}
	destportNameTxt := listaval[len(listaval)-1]
	dportSplit1 := strings.Split(destportNameTxt, ",")
	if len(dportSplit1) > 1 {
		destportNameTxt = dportSplit1[0]
		// ,左边的字母，不要数字
		var result []string
		for _, r := range destportNameTxt {
			if unicode.IsLetter(r) || unicode.IsSpace(r) {
				result = append(result, string(r))
			}
		}
		destportNameTxt = strings.Join(result, "")
		info.DestPortName = strings.TrimSpace(destportNameTxt)
	}

	fmt.Printf("------PoSheetDataParseA8aso--deliveryDateCustomerTxt(%s)--DeliveryDateCustomer(%+v)--info.DestPortName(%+v)-----\n", deliveryDateCustomerTxt, info.DeliveryDateCustomer, info.DestPortName)
	styleNoColIndex := -1
	sizeColIndex := 0
	qtyColIndex := 0
	for _, row := range rows {
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		// rowindex = uint(i + 1)
		// 跳出空数据行
		vala := strings.TrimSpace(row[0])
		// fmt.Printf("----PoSheetDataParseA8aso--eachrow(%+v)--styleNoColIndex(%d)-vala(%s)---\n", row, styleNoColIndex, vala)
		if vala == "" {
			continue
		}
		if styleNoColIndex == -1 {
			for coli, rv := range row {
				rv = strings.TrimSpace(rv)
				// 定位款号列
				if rv == "No." {
					styleNoColIndex = coli
				}
				// 定位尺码列
				if rv == "Variant" {
					sizeColIndex = coli
				}
				// 定位订单数量列
				if rv == "Quantity" {
					qtyColIndex = coli
				}
			}
		} else {
			// 有效的订单明细行
			// 款号 取每个sheet页的A列数据，从A9开始
			styleNo := strings.TrimSpace(row[styleNoColIndex])
			// 英文颜色 取每个sheet页的F列数据"-"前的数字为颜色，从F9开始
			// 尺码  取每个sheet页的F列数据"-"后的字母为尺码，从F9开始
			colorEn := ""
			size := ""
			// sizetxt := getCellTrimSpace(f, sheetName, "F", rowindex)
			sizetxt := strings.TrimSpace(row[sizeColIndex])
			sizesplit := strings.Split(sizetxt, "-")
			if len(sizesplit) > 1 {
				colorEn = sizesplit[0]
				size = sizesplit[1]
			}
			// 订单数量*	取每个sheet页的K列数据，从K9开始
			qtyStr := GetDigits(row[qtyColIndex])
			qty, err = strconv.Atoi(qtyStr)
			if err != nil {
				continue
			}

			item := OrderItem{}
			item.PoNo = info.PoNo
			item.StyleNo = styleNo                                      // 客户款号
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

	}
	return nil
}
