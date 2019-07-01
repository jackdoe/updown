package main

import (
	"fmt"
	"os"

	"io/ioutil"

	"github.com/yosssi/gohtml"
)

func main() {
	v, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	fmt.Println(gohtml.Format(string(v)))
}
