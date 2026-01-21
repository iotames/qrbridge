package biz

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA63amTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, 1, PoSheetDataParseA63am)
}

// 从Excel的单个sheet页面解析数据
func PoSheetDataParseA63am(f *excelize.File, sheetIndex int, info *PoInfo) error {

	var rows [][]string
	var err error
	var rowindex uint
	var qty int

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	if sheetName != "Summary" {
		// 目的国 取C6最后一个","右侧的英文
		dCountryText := getCellTrimSpace(f, sheetName, "C", 6)
		dCountryTextSplit := strings.Split(dCountryText, ",")
		info.DestCountry = strings.TrimSpace(dCountryTextSplit[len(dCountryTextSplit)-1])

		// 2025/12/21 or 12-21-25 D1
		deliveryDateCustomerText := getCellTrimSpace(f, sheetName, "D", 1)
		// 客户交期
		info.DeliveryDateCustomer, err = time.Parse("2006/01/02", deliveryDateCustomerText)
		if err != nil {
			// 12-21-25   "01-02-06"对应：月-日-缩写年（2位）
			info.DeliveryDateCustomer, _ = time.Parse("01-02-06", deliveryDateCustomerText)
		}

		info.PoNo = sheetName
		fmt.Printf("----PoNo(%s)----deliveryDateCustomerText(%s)--deliveryDateCustomerStr(%s)---\n", info.PoNo, deliveryDateCustomerText, info.DeliveryDateCustomer)
	}
	fmt.Printf("------PoSheetDataParseA63am-----info.DestCountry(%+v)--info.DestPortName(%+v)-----\n", info.DestCountry, info.DestPortName)

	if !info.DeliveryDateCustomer.IsZero() {
		// 客户交期
		deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

		// 离厂交期
		deliveryDateFactoryLeave := info.DeliveryDateCustomer.AddDate(0, 0, -7) // 离厂交期=客户交期-7
		deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")

		// 工厂交期
		deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, -7) // 工厂交期=离厂交期-7天
		deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
	}

	rows, err = f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("获取%s总行数失败: %w", sheetName, err)
	}

	for i, row := range rows {
		// fmt.Printf("----PoSheetDataParseA63am--eachrow(%+v)---\n", row)
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// 跳出空数据行
		if strings.TrimSpace(row[0]) == "" {
			continue
		}
		if rowindex < 10 {
			// 跳过前9行
			continue
		}
		// 客户款号。B10 开始
		styleNoText := getCellTrimSpace(f, sheetName, "A", rowindex)
		styleNoSplit := strings.Split(styleNoText, "-")
		lenStyleNoSplit := len(styleNoSplit)
		if lenStyleNoSplit < 2 {
			continue
		}
		styleNo := styleNoSplit[0]
		// 尺码
		size := styleNoSplit[lenStyleNoSplit-1]

		// 英文颜色。E11开始
		colorEn := getCellTrimSpace(f, sheetName, "F", rowindex)

		// 订单数量。J10开始
		qtytext := getCellTrimSpace(f, sheetName, "C", rowindex)
		if styleNo == "" || colorEn == "" || size == "" || qtytext == "" {
			continue
		}
		qtystr := GetDigits(qtytext)
		qty, err = strconv.Atoi(qtystr)
		if err != nil {
			continue
		}
		// 目的港
		destPortName := getCellTrimSpace(f, sheetName, "D", rowindex)

		item := OrderItem{}
		item.PoNo = info.PoNo
		item.StyleNo = styleNo
		item.Qty = qty                                              // 订单数量
		item.Size = size                                            // 尺码
		item.ColorEn = colorEn                                      // 英文颜色
		item.DestCountry = info.DestCountry                         // 目的国
		item.DestPortName = destPortName                            // 目的港
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		info.OrderItems = append(info.OrderItems, item)
	}
	return nil
}
