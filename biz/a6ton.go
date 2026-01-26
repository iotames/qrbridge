package biz

import (
	"fmt"
	// "strconv"
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
	// var qty int

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	if sheetName == "Summary" {
		// 客户交期：取“Summary”sheet页“Date Required”右侧的日期，再减去60天
		// 离厂交期=客户交期-15天
		// 工厂交期=离厂交期

		info.DestCountry = "Ireland"
		info.DestPortName = "Dublin"

		// 2025/12/21 or 12-21-25 D1
		deliveryDateCustomerText := getCellTrimSpace(f, sheetName, "D", 1)
		// 客户交期
		info.DeliveryDateCustomer, err = time.Parse("Jan-06", deliveryDateCustomerText)
		if err != nil {
			info.DeliveryDateCustomer, _ = time.Parse("2006/1/2", deliveryDateCustomerText)
		}
		// PO号 取“Summary”sheet页里，“PO: ”右侧的内容
		poText := getCellTrimSpace(f, sheetName, "A", 6)
		info.PoNo = strings.TrimSpace(strings.Replace(poText, "PO:", "", 1))

		// for i, row := range rows {
		// 	// fmt.Printf("----PoSheetDataParseA63am--eachrow(%+v)---\n", row)
		// 	// 当前行没有任何数据。跳过。
		// 	if len(row) == 0 {
		// 		continue
		// 	}
		// 	// 定义当前行号
		// 	rowindex = uint(i + 1)
		// 	// 跳出空数据行
		// 	cellA := strings.TrimSpace(row[0])
		// 	if cellA == "" {
		// 		continue
		// 	}
		// 	if strings.HasPrefix(cellA, "PO:") {
		// 		// 取“Summary”sheet页里，“PO: ”右侧的内容
		// 		info.PoNo = strings.TrimSpace(strings.Replace(cellA, "PO:", "", 1))
		// 	}
		// }

		for i, item := range info.OrderItems {
			item.PoNo = info.PoNo
			// item.ColorEn = colorEn                                      // 英文颜色
			item.DestCountry = info.DestCountry                         // 目的国
			item.DestPortName = info.DestPortName                       // 目的港
			item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
			item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
			item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。
			info.OrderItems[i] = item
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

		fmt.Printf("----PoNo(%s)----deliveryDateCustomerText(%s)--deliveryDateCustomerStr(%s)---\n", info.PoNo, deliveryDateCustomerText, info.DeliveryDateCustomer)
		return nil
	}
	fmt.Printf("------PoSheetDataParseA63am-----info.DestCountry(%+v)--info.DestPortName(%+v)-----\n", info.DestCountry, info.DestPortName)

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}

	// TODO 订单数量
	// 1、根据客户款号第二部分的内容（示例：2S、3S）去匹配sheet页，如果不是sheet名不是2S、3S，则取最后一个sheet页
	// 2、在匹配的sheet页里，通过“Colour/Print”列的内容（示例：ANT - AMB/MNE/WHI）去匹配款色，再通过尺码从款色下面提取对应的数量

	if sheetName == "SKU's" {
		// shortDescIndex := 0
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
			lastSyleNo := styleNoSplit[len(styleNoText)-1]
			styleNo := ""
			if lastSyleNo == "2S" || lastSyleNo == "3S" {
				styleNo = fmt.Sprintf("%s-%s", firstStyleNo, lastSyleNo)
			} else {
				styleNo = firstStyleNo
			}
			// 英文颜色。E2开始。
			// 取“SKU's”sheet页的“Colour/Print”列，“-”右侧的所有内容
			colorEnText := getCellTrimSpace(f, sheetName, "E", rowindex)
			if colorEnText == "Colour/Print" {
				continue
			}
			colorEnSplit := strings.Split(colorEnText, "-")
			colorEn := ""
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
			descSplit := strings.Split(desc, " ")
			lenDescSplit := len(descSplit)
			if lenDescSplit < 3 {
				continue
			}
			// 取倒数3个元素的切片
			descSplit = descSplit[lenDescSplit-3:]
			hasAge := false
			sizeSplit := []string{}
			for _, v := range descSplit {
				v := strings.TrimSpace(v)
				if hasAge {
					sizeSplit = append(sizeSplit, v)
				}
				if v == "Age" {
					hasAge = true
					sizeSplit = append(sizeSplit, v)
				}
			}
			size := ""
			if hasAge {
				size = strings.Join(sizeSplit, " ")
			} else {
				size = getCellTrimSpace(f, sheetName, "F", rowindex)
				size = strings.Replace(size, `'`, ``, 1)
			}

			item := OrderItem{}
			item.PoNo = info.PoNo
			item.StyleNo = styleNo
			item.ColorEn = colorEn
			item.Size = size
			info.OrderItems = append(info.OrderItems, item)
		}
		return nil
	}

	// for i, row := range rows {
	// 	// fmt.Printf("----PoSheetDataParseA63am--eachrow(%+v)---\n", row)
	// 	// 当前行没有任何数据。跳过。
	// 	if len(row) == 0 {
	// 		continue
	// 	}
	// 	// 定义当前行号
	// 	rowindex = uint(i + 1)
	// 	// 跳出空数据行
	// 	cellA := strings.TrimSpace(row[0])
	// 	if cellA == "" {
	// 		continue
	// 	}

	// 	// if rowindex < 10 {
	// 	// 	// 跳过前9行
	// 	// 	continue
	// 	// }
	// 	// qtystr := GetDigits(qtytext)
	// 	// qty, err = strconv.Atoi(qtystr)
	// 	// if err != nil {
	// 	// 	continue
	// 	// }
	// 	// 目的港

	// }
	return nil
}
