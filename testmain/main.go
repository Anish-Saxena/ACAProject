package main

import (
	"flag"
	"hfs"
	"log"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Please provide a mountpoint location")
	}
	hfs.BeginServer(flag.Arg(0))

}
