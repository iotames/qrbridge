package main

import (
	"flag"
)

var Debug, Dbinit, Daemon bool

func parseArgs() {
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.BoolVar(&Dbinit, "init", false, "init db")
	flag.BoolVar(&Daemon, "d", false, "run as daemon")
	flag.Parse()
}
