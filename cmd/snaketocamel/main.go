package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Printf("%s\n", strcase.ToCamel(s.Text()))
	}
}
