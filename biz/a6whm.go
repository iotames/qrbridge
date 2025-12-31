package biz

import (
	// "archive/zip"
	// "encoding/xml"
	// "bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func PoA6whmTransform(inputfile, outputfile string) (info PoInfo, err error) {
	return potransform(inputfile, outputfile, -1, PoSheetDataParseA6whm)
}

// 从Excel的每个sheet页面解析数据
func PoSheetDataParseA6whm(f *excelize.File, sheetIndex int, info *PoInfo) error {

	// 客户	    业务员	客户款号*	颜色*	英文颜色*	色号	PO NO*	尺码*	工厂交期	离厂交期*	客户交期*	订单数量*	目的国*	目的港	其他
	// A6WHM	张真真	有	无	有	无	有	有	离厂交期-7天	等于客户交期	有	有	有	无
	// 1、一个Excel对应一个PO，然后一个合同由多个PO组成

	var rows [][]string
	var err error
	// var poInt int
	var rowindex uint

	deliveryDateCustomerStr := ""     // 客户交期
	deliveryDateFactoryLeaveStr := "" // 离厂交期
	deliveryDateFactoryStr := ""      // 工厂交期
	sheetName := f.GetSheetName(sheetIndex)

	if sheetIndex == 0 {
		// 目的国
		info.DestCountry = getCellTrimSpace(f, sheetName, "E", 15) // THE NETHERLANDS
		info.PoNo = getCellTrimSpace(f, sheetName, "J", 11)
		deliveryDateCustomerTxt := getCellTrimSpace(f, sheetName, "M", 11)
		// 客户交期 日.月.年
		info.DeliveryDateCustomer, _ = time.Parse("02.01.2006", deliveryDateCustomerTxt)
	}

	if !info.DeliveryDateCustomer.IsZero() {
		// 客户交期
		deliveryDateCustomerStr = info.DeliveryDateCustomer.Format("2006-01-02")

		// 离厂交期
		deliveryDateFactoryLeave := info.DeliveryDateCustomer // 离厂交期==客户交期
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
		// 当前行没有任何数据。跳过。
		if len(row) == 0 {
			continue
		}
		// 定义当前行号
		rowindex = uint(i + 1)
		// 获取当前行的第一列数据
		rowA := strings.TrimSpace(row[0])

		// 提取客户款号
		// isOrderDetailRow
		orderId, err := strconv.Atoi(rowA)
		if err != nil {
			continue
		}
		if orderId < 30000000 {
			continue
		}
		// fmt.Printf("-----PoSheetDataParseA6whm--row(%+v)---\n", row)

		// 商品标题或者详情 在D列或B列
		desc := getCellTrimSpace(f, sheetName, "D", rowindex)
		colorEn, size := getA6whmColorSizeByDesc(desc)
		if colorEn == "" || size == "" {
			desc = getCellTrimSpace(f, sheetName, "B", rowindex)
			colorEn, size = getA6whmColorSizeByDesc(desc)
		}

		// 提取订单数量 在C列或H列
		// 一定要先检查C列，没有则取H列。否则可能直接取到H列的总金额。
		qtytext := getCellTrimSpace(f, sheetName, "C", rowindex)
		if qtytext == "" || qtytext == rowA {
			qtytext = getCellTrimSpace(f, sheetName, "H", rowindex)
		}
		qtystr := GetDigits(qtytext)
		qty, err := strconv.Atoi(qtystr)
		if err != nil {
			continue
		}

		// for coli, celval := range row {}
		// TODO "目的港"
		item := OrderItem{}
		item.PoNo = info.PoNo
		item.StyleNo = fmt.Sprintf("%d", orderId) // 客户款号
		item.Qty = qty                            // 订单数量
		item.Desc = desc
		item.Size = size       // 尺码
		item.ColorEn = colorEn // 英文颜色

		item.DestCountry = info.DestCountry                         // 目的国
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		info.OrderItems = append(info.OrderItems, item)
	}
	return nil
}

func getA6whmColorSizeByDesc(desc string) (color, size string) {
	descsplit := strings.Split(desc, ",")
	if len(descsplit) < 3 {
		return
	}
	return strings.TrimSpace(descsplit[2]), strings.TrimSpace(descsplit[1])
}
