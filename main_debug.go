package main

import (
	"fmt"

	"github.com/iotames/qrbridge/biz"

	// "github.com/iotames/qrbridge/db"
	// "github.com/iotames/qrbridge/hotswap"
	"github.com/iotames/qrbridge/service"
)

func debug() {
	// sqlTxt, err := db.GetSQL("pricing_percent.sql", "and cp.customer_name in(?, ?, ?)")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("--debug--sqlTxt----", sqlTxt)
	// filelist := hotswap.GetScriptDir(nil).LsDirByEmbedFS()
	// for _, f := range filelist {
	// 	fmt.Println(f)
	// }
	if Inputfile != "" {
		f, err := service.NewTableFile(Inputfile).OpenExcel()
		if err != nil {
			panic(fmt.Errorf("打开Excel文件失败: %w", err))
		}
		info := biz.PoInfo{}
		err = biz.PoSheetDataParseBewcw(f, 0, &info)
		if err != nil {
			panic(fmt.Errorf("解析Excel文件失败: %w", err))
		}
		for i, item := range info.OrderItems {
			fmt.Printf("----i(%d)---item(%+v)---------\n", i, item)
		}
	}
}
