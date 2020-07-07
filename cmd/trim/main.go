package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jackdoe/updown/util"
)

func main() {
	cut := flag.String("cut", "", "trim cutset")
	flag.Parse()

	util.ForeachLine(os.Stdin, func(line string, _hasNewLine bool) {
		var trimmed string
		if *cut != "" {
			trimmed = strings.Trim(line, *cut)
		} else {
			trimmed = strings.TrimSpace(line)
		}
		fmt.Println(trimmed)
	})
}
