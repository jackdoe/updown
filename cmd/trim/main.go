package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	cut := flag.String("cut", "", "trim cutset what?")
	flag.Parse()

	strings.TrimSpace("a")
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if *cut != "" {
			fmt.Println(strings.Trim(line[:len(line)-1], *cut))
		} else {
			fmt.Println(strings.TrimSpace(line[:len(line)-1]))
		}
	}
}
