package main

import (
	"aca/linkedfs"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Please enter 2 arguments: linked and original directory, respectively. Exiting now.")
	}
	link := flag.Arg(0)
	orig := flag.Arg(1)
	linkedfs.Begin(orig, link)
}
