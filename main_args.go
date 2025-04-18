package main

import (
	"flag"
)

var Dbinit, Daemon bool

func parseArgs() {
	flag.BoolVar(&Dbinit, "init", false, "init db")
	flag.BoolVar(&Daemon, "d", false, "run as daemon")
	flag.Parse()
}
