package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fn := flag.String("file", "-", "filename or - for stdin")
	flag.Parse()
	var content []byte
	var err error
	if *fn == "-" {
		content, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else {
		content, err = ioutil.ReadFile(*fn)
		if err != nil {
			panic(err)
		}
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, *fn, string(content), parser.ParseComments)
	if err == nil {
		for _, s := range file.Imports {
			fmt.Printf("%s %s\n", file.Name, strings.Trim(s.Path.Value, `"`))
		}
	}
}
