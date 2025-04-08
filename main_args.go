package main

import (
	"flag"
)

var Dbinit bool

func parseArgs() {
	flag.BoolVar(&Dbinit, "init", false, "init db")
	flag.Parse()
}
