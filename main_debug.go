package main

import (
	"fmt"

	"github.com/iotames/qrbridge/sql"
)

func debug() {
	for _, f := range sql.LsDir() {
		fmt.Println(f)
	}
}
