package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jackdoe/updown/util"
)

func main() {
	bucket := flag.Int("bucket", 100, "bucket size")
	flag.Parse()
	b := int64(*bucket)
	util.ForeachLine(os.Stdin, func(text string, _last bool) {
		v := int64(util.FloatOrPanic(text))
		fmt.Printf("%d\n", (v/b)*b)
	})
}
