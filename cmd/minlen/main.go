package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jackdoe/updown/util"
)

func main() {
	pn := flag.Int("n", 1, "min line length")
	flag.Parse()

	util.ForeachLine(os.Stdin, func(lineNL string, end bool) {
		n := *pn
		if !end {
			n-- // ignore \n
		}
		if len(lineNL) >= n {
			fmt.Print(lineNL)
		}
	})
}
