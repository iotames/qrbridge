package main

import (
	"fmt"

	"github.com/iotames/qrbridge/sql"
)

func debug() {
	sqlTxt, err := sql.GetSQL("pricing_percent.sql", "and cp.customer_name in(?, ?, ?)")
	if err != nil {
		panic(err)
	}
	fmt.Printf("-----------sqlTxt(%s)--------", sqlTxt)
}
