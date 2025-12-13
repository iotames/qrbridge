package main

import (
	"flag"
)

var Debug, DbPing, Daemon, vsion bool

func parseArgs() {
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.BoolVar(&DbPing, "dbping", false, "ping db")
	flag.BoolVar(&Daemon, "d", false, "run as daemon")
	flag.BoolVar(&vsion, "version", false, "show version")
	flag.Parse()
}
