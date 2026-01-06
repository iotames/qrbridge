package biz

import (
	"fmt"

	"github.com/iotames/qrbridge/service"
	"github.com/xuri/excelize/v2"
)

var poCommonTplTitleRow = []interface{}{"客户款号*", "颜色*", "英文颜色*", "色号", "PO NO*", "尺码*", "工厂交期", "离厂交期*", "客户交期*", "订单数量*", "目的国*", "目的港"}

func poOutputExcel(outputfile string, info PoInfo) error {
	var f *excelize.File
	var err error
	// 输出新的EXCEL
	f = service.NewTableFile(outputfile).NewExcel()

	titleRow := poCommonTplTitleRow
	err = f.SetSheetRow("Sheet1", "A1", &titleRow)
	if err != nil {
		return fmt.Errorf("fai to SetSheetRow: %w", err)
	}
	for i := 0; i < len(info.OrderItems); i++ {
		rowIndex := i + 2
		err = f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", rowIndex), &[]interface{}{info.OrderItems[i].StyleNo, info.OrderItems[i].Color, info.OrderItems[i].ColorEn, info.OrderItems[i].ColorNo, info.OrderItems[i].PoNo, info.OrderItems[i].Size, info.OrderItems[i].DeliveryDateFactory, info.OrderItems[i].DeliveryDateFactoryLeave, info.OrderItems[i].DeliveryDateCustomer, info.OrderItems[i].Qty, info.OrderItems[i].DestCountry, info.OrderItems[i].DestPortName})
		if err != nil {
			return fmt.Errorf("fai to SetSheetRow[%d]: %w", rowIndex, err)
		}
	}
	err = f.SaveAs(outputfile)
	if err != nil {
		return fmt.Errorf("保存%s文件失败: %w", outputfile, err)
	}
	err = f.Close()
	return err
}
