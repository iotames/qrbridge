package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA3mjmTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseA3mjm)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseA3mjm(f *excelize.File, sheetIndex int, info *PoInfo) error {
	var rows [][]string
	var err error
	var rowindex uint
	var qty int
	// var desc, colorEn, size string

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}

	// 目的港
	info.DestPortName = "VALENCIA"

	startOrderData := false
	var sizeList []string
	for i, row := range rows {
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// 跳出空数据行
		vala := strings.TrimSpace(row[0])
		// fmt.Printf("----PoSheetDataParseA3mjm--eachrow(%+v)--styleNoColIndex(%d)-vala(%s)---\n", row, styleNoColIndex, vala)
		if vala == "" {
			continue
		}
		// 4. 固定取第一个sheet页“P/O No:”右侧的数字
		if strings.Contains(vala, "P/O No:") {
			poinfoSplit := strings.Split(vala, "\n")
			fmt.Printf("-----poinfoSplit(%s)---\n", strings.Join(poinfoSplit, "|"))
			if len(poinfoSplit) > 0 {
				info.PoNo = strings.TrimSpace(strings.Replace(poinfoSplit[0], "P/O No:", "", 1))
				for _, pocell := range poinfoSplit {
					// DESTINATION:     JOMA SPAIN
					if strings.Contains(pocell, "DESTINATION:") {
						// 目的地	取“DESTINATION: ”右侧的单词，去掉“JOMA”
						replacer := strings.NewReplacer("DESTINATION:", "", "JOMA", "")
						info.DestCountry = strings.TrimSpace(replacer.Replace(pocell))
					}
					if strings.Contains(pocell, "EX-WORKS DATE:") {
						// 客户交期	取第一个sheet页“EX-WORKS DATE:”右侧的日期
						// "Mar 28 th, 26"
						deliveryDateCustomerTxt := strings.TrimSpace(strings.Replace(pocell, "EX-WORKS DATE:", "", 1))

						// 1. 移除常见的序数后缀
						// 创建一个副本进行操作，避免修改原字符串
						dateToParse := deliveryDateCustomerTxt
						replacer := strings.NewReplacer(" th,", ",", " st,", ",", " nd,", ",", " rd,", ",")
						dateToParse = replacer.Replace(dateToParse)
						// 清理可能多余的空格
						dateToParse = strings.ReplaceAll(dateToParse, "  ", " ")
						dateToParse = strings.TrimSpace(dateToParse)

						// 2. 定义布局字符串
						// 现在 dateToParse 应该是 "Mar 28, 26"
						// Jan (月份缩写) -> 对应 "Mar"
						// 02 (两位数的日期) -> 对应 "28"
						// 06 (两位数的年份) -> 对应 "26"

						info.DeliveryDateCustomer, _ = time.Parse("Jan 02, 06", dateToParse)
						fmt.Printf("-----deliveryDateCustomerTxt(%s)--dateToParse(%s)--DeliveryDateCustomer(%+v)--\n", deliveryDateCustomerTxt, dateToParse, info.DeliveryDateCustomer)
					}
				}

			}

			continue
		}

		if vala == "Style" {
			startOrderData = true
			// 5. 尺码 取每个sheet页“Tariff”右侧的内容，截止到“Total”前面
			for _, titleCell := range row {
				if strings.Contains(titleCell, "Tariff") {

					sizeSplit := strings.Split(titleCell, " ")
					for _, sz := range sizeSplit {
						sz = strings.TrimSpace(sz)
						if sz == "Tariff" || sz == "Total" || sz == "" {
							continue
						}
						sizeList = append(sizeList, sz)
					}
					fmt.Printf("---rowTitle(%s)---sizeList(%s)----\n", titleCell, strings.Join(sizeList, "|"))
				}
			}

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
			fmt.Printf("------PoSheetDataParseA3mjm----DeliveryDateCustomer(%+v)--info.DestPortName(%+v)-----\n", info.DeliveryDateCustomer, info.DestPortName)
		}
		if startOrderData {
			// 有效的订单明细行
			rowSplit := strings.Split(vala, "  ")
			rowDataList := []string{}
			for _, rowd := range rowSplit {
				rowd = strings.TrimSpace(rowd)
				if rowd == "" {
					continue
				}
				rowDataList = append(rowDataList, rowd)
			}
			fmt.Printf("----startOrderData--rowindex(%d)--rowDataList(%s)---\n", rowindex, strings.Join(rowDataList, "|"))
			if len(rowSplit) < 5 {
				continue
			}
			// "902400.346         R-NATURE SHORT GREEN TURQUOISE   62.03.43.90      130      210       160        70         20                                                                                                                                                                                                               590"

			// 1. 款号 取每个sheet页A列“Style”下方表格最左侧第一个空格前的内容
			styleNo := strings.TrimSpace(rowDataList[0])

			// 2. 英文颜色 取每个sheet页A列“Style”下方表格内的英文内容，需要业务员在标准模板中手动处理
			colorEn := strings.TrimSpace(rowDataList[1])

			// 3. 色号 取款号固定后三位数数字，不含小数点
			colorNo := ""
			styleNoSplit := strings.Split(styleNo, ".")
			if len(styleNoSplit) > 1 {
				colorNo = styleNoSplit[len(styleNoSplit)-1]
			}

			// 5. 尺码 取每个sheet页“Tariff”右侧的内容，截止到“Total”前面
			sizeQtyList := rowDataList[3 : len(rowDataList)-1]
			fmt.Printf("-----sizeList(%s)--sizeQty(%s)-------\n", strings.Join(sizeList, "|"), strings.Join(sizeQtyList, "|"))
			if len(sizeQtyList) == len(sizeList) {

				for iii, qtystr := range sizeQtyList {
					qtystr = strings.TrimSpace(qtystr)
					size := sizeList[iii]
					qtystr = GetDigits(qtystr)
					qty, _ = strconv.Atoi(qtystr)

					item := OrderItem{}
					item.PoNo = info.PoNo
					item.StyleNo = styleNo                                      // 客户款号
					item.ColorNo = colorNo                                      // 色号
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

		}

	}
	return nil
}
