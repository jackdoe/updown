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

	util.ForeachLine(os.Stdin, func(line string, _hasNewLine bool) {
		if len(line) >= *pn {
			fmt.Println(line)
		}
	})
}
