package biz

import (
	// "bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	// go get github.com/unidoc/unipdf/v4
	// "github.com/ledongthuc/pdf"
	"github.com/xuri/excelize/v2"
)

func PoBewcwTransform(inputfile, outputfile string) (info PoInfo, err error) {
	// pdf.DebugOn = true
	// var content string
	// content, err = readPdf(inputfile) // Read local pdf file
	// var ff *os.File
	// ff, _ = os.OpenFile(inputfile, os.O_CREATE|os.O_TRUNC, 0755)
	// defer ff.Close()
	// _, err = io.WriteString(ff, content)
	// fmt.Println(content)
	return potransform(inputfile, outputfile, 0, PoSheetDataParseBewcw)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseBewcw(f *excelize.File, sheetName string, info *PoInfo) error {
	var rows [][]string
	var err error
	var poInt int
	var rowindex uint
	var deliveryDateCustomer time.Time
	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}
	for i, row := range rows {
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// 获取当前行的第一列数据rowA
		rowA := strings.TrimSpace(row[0])
		// fmt.Printf("----rowindex(%d)--rowA(%s)--\n", rowindex, rowA)

		// 当前rowA数据是PO编号，设置正整数的PO号，然后继续解析下一行
		if strings.HasPrefix(rowA, "No.") {
			if strings.Contains(rowA, "PO") {
				// 字符串是No.开头且包含PO字样的单元格，确认为PO数据，开始解析
				// 1. 去掉No.字符串，去掉PO字符串，去掉首尾空白符。
				// 2. 转为int整数
				postr := strings.TrimSpace(strings.Replace(strings.Replace(rowA, "No.", "", 1), "PO", "", 1))
				// 设置PO号
				poInt, err = strconv.Atoi(postr)
				if err == nil {
					// PO号提取成功
					continue
				}
			}
		}

		// 当前rowA数据包含Shipment Date
		if strings.Contains(rowA, "Shipment Date") {
			if len(row) > 2 {
				dateText := row[2]
				datesplit := strings.Split(dateText, "\n")
				// fmt.Printf("------row(%+v)---dateText(%s)----datesplit(%+v)----\n", row, dateText, datesplit)
				// for ii, rcc := range row {
				// 	fmt.Printf("---rowindex(%d)--ii(%d)-----rcc(%s)----\n", rowindex, ii, rcc)
				// }
				if len(datesplit) > 1 {
					deliveryDateCustomerStr = "20" + datesplit[1] // 客户交期

					deliveryDateCustomer, err = time.Parse("2006-01-02", deliveryDateCustomerStr)
					if err == nil {
						// 客户交期的前7天是离厂交期，离厂交期的前5天是工厂交期。
						deliveryDateFactoryLeave := deliveryDateCustomer.AddDate(0, 0, -7) // 离厂交期。必填。客户交期-7天
						deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7)  // 工厂交期。非必填。离厂交期-5天
						deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")
						deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
					}
				}
			}
		}

		// 跳过订单详情区块的标题行
		if rowA == "No." {
			continue
		}

		// 提取客户款号
		// isOrderDetailRow
		orderId, err := strconv.Atoi(rowA)
		if err != nil {
			continue
		}
		if orderId < 10000 {
			continue
		}

		// 提取订单数量
		qtytext := getCellTrimSpace(f, sheetName, "I", rowindex)
		qty, err := strconv.Atoi(strings.Replace(qtytext, ",", "", 1))
		if err != nil {
			continue
		}
		desc := getCellTrimSpace(f, sheetName, "B", rowindex)

		// 提取尺码
		sizeText := getCellTrimSpace(f, sheetName, "F", rowindex)
		sizesplit := strings.Split(sizeText, "-")
		size := sizeText

		// 提取颜色号
		colorNo := ""
		if len(sizesplit) > 1 {
			size = sizesplit[1]
			colorNo = sizesplit[0]
		}

		// 提取英文颜色
		colorEn := ""
		// 端横杠-分隔取第二端字符串
		descsplit := strings.Split(desc, "W-")
		if len(descsplit) > 1 {
			colorEn = strings.TrimSpace(strings.Replace(descsplit[1], "-"+size, "", 1))
		}
		// for coli, celval := range row {}
		// TODO "目的港"
		item := OrderItem{}
		item.PoNo = fmt.Sprintf("PO%d", poInt)
		item.StyleNo = fmt.Sprintf("%d", orderId) // 客户款号
		item.Qty = qty                            // 订单数量
		item.Desc = desc
		item.Size = size                                            // 尺码
		item.ColorEn = colorEn                                      // 英文颜色
		item.ColorNo = colorNo                                      // 颜色号
		item.DestCountry = "Sweden"                                 // 目的国
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。客户交期-7天
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-5天
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
