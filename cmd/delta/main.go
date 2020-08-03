package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jackdoe/updown/util"
)

func main() {
	isFloat := flag.Bool("f", false, "expect floating point numbers")
	isInverse := flag.Bool("i", false, "inverse")
	flag.Parse()

	if !*isFloat {
		prev := int64(0)
		util.ForeachLine(os.Stdin, func(text string, _last bool) {
			v := util.IntOrPanic(text)
			if !*isInverse {
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
			if !*isInverse {
				fmt.Println(v - prev)
			} else {
				fmt.Println(prev - v)
			}
			prev = v
		})
	}
}
