package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jackdoe/updown/util"
)

func main() {
	isInt := flag.Bool("int", false, "expect int numbers only")
	desc := flag.Bool("desc", false, "expect descending order")
	flag.Parse()

	if *isInt {
		prev := int64(0)
		util.ForeachLine(os.Stdin, func(text string, _last bool) {
			v := util.IntOrPanic(text)
			if !*desc {
				fmt.Println(v - prev)
			} else {
				fmt.Println(prev - v)
			}
			prev = v
		})
	} else {
		prev := float64(0)
		util.ForeachLine(os.Stdin, func(text string, _last bool) {
			v := util.FloatOrPanic(text)
			if !*desc {
				fmt.Println(v - prev)
			} else {
				fmt.Println(prev - v)
			}
			prev = v
		})
	}
}
