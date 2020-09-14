package main

import (
	"os"

	"github.com/iancoleman/strcase"
	"github.com/jackdoe/updown/util"
)

func main() {
	util.ForeachLine(os.Stdin, func(text string, hasNewLine bool) {
		os.Stdout.Write([]byte(strcase.ToSnake(text)))
		if hasNewLine {
			os.Stdout.Write(util.NL)
		}
	})
}
