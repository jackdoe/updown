package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jackdoe/updown/util"
)

func main() {
	isInt := flag.Bool("int", false, "expect int numbers only")
	flag.Parse()

	if *isInt {
		sum := int64(0)
		util.ForeachLine(os.Stdin, func(text string, _last bool) {
			v := util.IntOrPanic(text)
			sum += v
		})
		fmt.Printf("%d\n", sum)

	} else {
		sum := float64(0)
		util.ForeachLine(os.Stdin, func(text string, _last bool) {
			v := util.FloatOrPanic(text)
			sum += v
		})
		fmt.Printf("%f\n", sum)
	}
}
