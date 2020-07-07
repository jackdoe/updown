package main

import (
	"fmt"
	"os"

	"github.com/jackdoe/updown/util"
)

func main() {
	uniq := map[string]uint64{}
	util.ForeachLine(os.Stdin, func(line string, _hasNewLine bool) {
		uniq[line]++
	})

	for k, count := range uniq {
		fmt.Printf("%8d %s\n", count, k)
	}
}
