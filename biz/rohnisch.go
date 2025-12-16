package biz

import (
	"fmt"
	"strings"
	"time"

	"github.com/iotames/qrbridge/service"
	"github.com/xuri/excelize/v2"
)

func PoRohnischTransform(inputtpl, inputfile, outputfile string) (info PoInfo, err error) {
	f, err := service.NewTableFile(inputfile).OpenExcel()
	if err != nil {
		return PoInfo{}, fmt.Errorf("打开Excel文件失败: %w", err)
	}

	sheets := f.GetSheetList()
	for i, sheet := range sheets {
		poSheetDataParseRohnisch(f, sheet, i, &info)
	}

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
func poSheetDataParseRohnisch(f *excelize.File, sheetName string, id int, info *PoInfo) error {
	// 第57行为标题行
	// 客户款号C58
	// 款式标题G58 Kay High Waist Tights  Black/Black, XS,  0101
	// noTitle, _ := f.GetCellValue(sheetName, "C57")
	// no, _ := f.GetCellValue(sheetName, "C58")
	poNo := getCellTrimSpace(f, sheetName, "T", 10)                    // 获取PO号
	deliveryDateCustomerTxt := getCellTrimSpace(f, sheetName, "I", 50) // 获取客户交期。I50 25-09-15
	deliveryDateCustomerStr := "20" + deliveryDateCustomerTxt          // 客户交期
	deliveryDateFactoryLeaveStr := ""                                  // 离厂交期
	deliveryDateFactoryStr := ""                                       // 工厂交期
	deliveryDateCustomer, err := time.Parse("2006-01-02", deliveryDateCustomerStr)
	if err == nil {
		deliveryDateFactoryLeave := deliveryDateCustomer.AddDate(0, 0, -7)
		deliveryDateFactory := deliveryDateFactoryLeave.AddDate(0, 0, 7)
		deliveryDateFactoryLeaveStr = deliveryDateFactoryLeave.Format("2006-01-02")
		deliveryDateFactoryStr = deliveryDateFactory.Format("2006-01-02")
	}
	destCountry := getCellTrimSpace(f, sheetName, "F", 31) // 目的国
	ids := getOkOrderItemRowIndexs(f, sheetName, "C", 58, 5, []string{"No."})
	for _, rowindex := range ids {
		item := OrderItem{}
		item.StyleNo = getCellTrimSpace(f, sheetName, "C", rowindex) // 提取客户款号
		item.Desc = getCellTrimSpace(f, sheetName, "G", rowindex)
		item.ColorEn, item.Size = getRohnischColorSizeByDesc(item.Desc) // 获取颜色尺码
		item.PoNo = poNo
		item.DeliveryDateCustomer = deliveryDateCustomerStr         // 客户交期。必填。
		item.DeliveryDateFactoryLeave = deliveryDateFactoryLeaveStr // 离厂交期。必填。客户交期-7天
		item.DeliveryDateFactory = deliveryDateFactoryStr           // 工厂交期。非必填。离厂交期-7天
		qtyStr := getCellTrimSpace(f, sheetName, "AA", rowindex)    // 原始数据。订单数量。必填
		fmt.Sscanf(qtyStr, "%d", &item.Qty)                         // 转换为整型。订单数量。必填
		item.DestCountry = destCountry
		info.OrderItems = append(info.OrderItems, item)
		fmt.Printf("----sheet(%d-%s)---rowindex(%d)---orderItem(%+v)------\n", id, sheetName, rowindex, item)
	}
	return nil
}

func getRohnischColorSizeByDesc(desc string) (color, size string) {
	descsplit := strings.Split(desc, ",")
	splitlen := len(descsplit)
	if splitlen > 2 {
		size = strings.TrimSpace(descsplit[splitlen-2]) // 倒数第二个
		titleColor := strings.TrimSpace(descsplit[0])
		colorsplit := strings.Split(titleColor, " ")
		colorsplitlen := len(colorsplit)
		if colorsplitlen > 1 {
			color = strings.TrimSpace(colorsplit[colorsplitlen-1])
		}
	}
	return
}
