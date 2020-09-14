package main

import (
	"flag"
	"os"

	"github.com/jackdoe/updown/util"
)

func main() {
	pskip := flag.Int("n", 1, "skip first N lines")
	flag.Parse()

	skip := *pskip
	util.ForeachLine(os.Stdin, func(text string, hasNewLine bool) {
		if skip > 0 {
			skip--
			return
		}

		os.Stdout.Write([]byte(text))
		if hasNewLine {
			os.Stdout.Write(util.NL)
		}
	})
}
