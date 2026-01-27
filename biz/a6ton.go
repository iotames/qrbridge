package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA6tonTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, poSheetDataParseA6ton)
}

// 从Excel的每个sheet页面解析数据
func poSheetDataParseA6ton(f *excelize.File, sheetIndex int, info *PoInfo) error {

	var rows [][]string
	var err error
	var rowindex uint
	var found bool
	// var qty int
	info.DestCountry = "Ireland"
	info.DestPortName = "Dublin"
	deliveryDateCustomerText := ""
	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)
	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}

	if sheetName == "Summary" {

		for i, row := range rows {
			// fmt.Printf("----PoSheetDataParseA63am--eachrow(%+v)---\n", row)
			// 当前行没有任何数据。跳过。
			if len(row) == 0 {
				continue
			}
			// 定义当前行号
			rowindex = uint(i + 1)
			// 跳出空数据行
			cellA := strings.TrimSpace(row[0])
			if cellA == "" {
				continue
			}

			// PO号 取“Summary”sheet页里，“PO: ”右侧的内容
			// poText := getCellTrimSpace(f, sheetName, "A", 6)
			if strings.HasPrefix(cellA, "PO:") {
				info.PoNo = strings.TrimSpace(strings.Replace(cellA, "PO:", "", 1))
				break
			}

			// 客户交期：取“Summary”sheet页“Date Required”右侧的日期，再减去60天
			// 离厂交期=客户交期-15天
			// 工厂交期=离厂交期
			// 2025/12/21 or 12-21-25 D1
			if info.DeliveryDateCustomer.IsZero() {
				deliveryDateCustomerText, found = GetNextStringTrimSpace(row, "Date Required")
				if found && deliveryDateCustomerText != "" {
					// 客户交期
					info.DeliveryDateCustomer, err = time.Parse("Jan-06", deliveryDateCustomerText)
					info.DeliveryDateCustomer = info.DeliveryDateCustomer.AddDate(0, -60, 0)
					if err != nil {
						info.DeliveryDateCustomer, _ = time.Parse("2006/1/2", deliveryDateCustomerText)
					}
				}
			}
		}

		if !info.DeliveryDateCustomer.IsZero() {
			// 客户交期
			deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

			// 离厂交期
			deliveryDateFactoryLeave := info.DeliveryDateCustomer.AddDate(0, 0, -15) // 离厂交期=客户交期-15
			deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")

			// 工厂交期
			deliveryDateFactory := deliveryDateFactoryLeave // 工厂交期=离厂交期
			deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
		}

		// 更新每个尺码的PO, 目的国，目的港，客户交期，离厂交期，工厂交期
		for i, item := range info.OrderItems {
			item.PoNo = info.PoNo
			item.DestCountry = info.DestCountry                         // 目的国
			item.DestPortName = info.DestPortName                       // 目的港
			item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
			item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
			item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。
			info.OrderItems[i] = item
		}

		fmt.Printf("----PoNo(%s)----deliveryDateCustomerText(%s)--deliveryDateCustomerStr(%s)---\n", info.PoNo, deliveryDateCustomerText, info.DeliveryDateCustomer)
		return nil
	}
	fmt.Printf("------PoSheetDataParseA63am-----info.DestCountry(%+v)--info.DestPortName(%+v)-----\n", info.DestCountry, info.DestPortName)

	if sheetName == "SKU's" {
		sheets := f.GetSheetList()
		lenSheets := len(sheets)
		// shortDescIndex := 0
		rows2S, _ := f.GetRows("2S")
		rows3S, _ := f.GetRows("3S")
		lastSheetRows, _ := f.GetRows(sheets[lenSheets-1])
		orderRowsMap := map[string][][]string{
			"2S": rows2S,
			"3S": rows3S,
		}
		var orderRows [][]string

		for i, row := range rows {
			// fmt.Printf("----PoSheetDataParseA63am--eachrow(%+v)---\n", row)
			// 当前行没有任何数据。跳过。
			if len(row) == 0 {
				continue
			}
			// 定义当前行号
			rowindex = uint(i + 1)
			// 跳出空数据行
			cellA := strings.TrimSpace(row[0])
			if cellA == "" {
				continue
			}
			// 客户款号。D2 开始
			// 1、第一部分取“SKU's”sheet页的“Short Description”列，截止到空格后的三位数字
			// 2、第二部分取“SKU's”sheet页的“Short Description”列的最后两个字符（只取2S和3S，如果没有则不取）
			// 3、将两个部分中间用“-”进行拼接

			styleNoText := getCellTrimSpace(f, sheetName, "D", rowindex)
			if styleNoText == "Short Description" {
				continue
			}
			styleNoSplit := strings.Split(styleNoText, " ")
			lenStyleNoSplit := len(styleNoSplit)
			if lenStyleNoSplit < 2 {
				continue
			}
			firstStyleNo := styleNoSplit[1]
			lastSyleNo := styleNoSplit[lenStyleNoSplit-1]
			styleNo := ""

			if lastSyleNo == "2S" || lastSyleNo == "3S" {
				styleNo = fmt.Sprintf("%s-%s", firstStyleNo, lastSyleNo)
				orderRows = orderRowsMap[lastSyleNo]
			} else {
				styleNo = firstStyleNo
				orderRows = lastSheetRows
			}

			// 英文颜色。E2开始。
			// 取“SKU's”sheet页的“Colour/Print”列，“-”右侧的所有内容
			colorEnText := getCellTrimSpace(f, sheetName, "E", rowindex)
			if colorEnText == "Colour/Print" {
				continue
			}
			colorEnSplit := strings.Split(colorEnText, "-")
			colorEn := colorEnText
			if len(colorEnSplit) > 1 {
				colorEn = strings.TrimSpace(colorEnSplit[1])
			}

			// 尺码

			// 先判断“Description”列“)”右侧的内容是否包含“Age”
			// 1、是：取“Description”列“)”右侧的所有内容
			// 2、否：取“Size”列的所有内容，需要去掉“'”号
			desc := getCellTrimSpace(f, sheetName, "C", rowindex)
			if desc == "Description" {
				continue
			}
			descSplit := strings.Split(desc, ") ")
			lenDescSplit := len(descSplit)
			if lenDescSplit < 2 {
				continue
			}
			// 原始的尺码列的内容
			sizeTxt := getCellTrimSpace(f, sheetName, "F", rowindex)
			size := ""
			if strings.Contains(descSplit[1], "Age") {
				size = strings.TrimSpace(descSplit[1])
			} else {
				size = strings.Replace(sizeTxt, `'`, ``, 1)
			}

			// 订单数量
			// 1、根据客户款号第二部分的内容（示例：2S、3S）去匹配sheet页，如果不是sheet名不是2S、3S，则取最后一个sheet页
			// 2、在匹配的sheet页里，通过“Colour/Print”列的内容（示例：ANT - AMB/MNE/WHI）去匹配款色，再通过尺码从款色下面提取对应的数量
			qty := getA6tonOrderQty(orderRows, colorEnText, sizeTxt)

			item := OrderItem{}
			item.PoNo = info.PoNo
			item.StyleNo = styleNo
			item.ColorEn = colorEn
			item.Size = size
			item.Qty = qty
			info.OrderItems = append(info.OrderItems, item)
		}
		return nil
	}
	return nil
}

func getA6tonOrderQty(orderRows [][]string, colorEnText string, sizeTxt string) int {
	var rowindex uint
	var colorRowIndex uint = 0
	for i, orderRow := range orderRows {
		// 当前行没有任何数据。跳过。
		if len(orderRow) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// 跳出空数据行
		orderCell1 := strings.TrimSpace(orderRow[0])
		if orderCell1 == "" {
			continue
		}
		key := orderCell1
		// 定位到颜色总数所在的行
		if key == colorEnText {
			colorRowIndex = rowindex
			continue
		}
		if colorRowIndex > 0 && rowindex > colorRowIndex {
			// 定位尺码所在行
			qtyTxt, found := GetNextStringTrimSpace(orderRow, sizeTxt)
			if found {
				qtystr := GetDigits(qtyTxt)
				qty, err := strconv.Atoi(qtystr)
				if err != nil {
					return 0
				}
				return qty
			}
		}
	}
	return 0
}
