package main

import (
	"fmt"

	"github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/hotswap"
)

func debug() {
	sqlTxt, err := db.GetSQL("pricing_percent.sql", "and cp.customer_name in(?, ?, ?)")
	if err != nil {
		panic(err)
	}
	fmt.Println("--debug--sqlTxt----", sqlTxt)
	filelist := hotswap.GetScriptDir(nil).LsDirByEmbedFS()
	for _, f := range filelist {
		fmt.Println(f)
	}
}
