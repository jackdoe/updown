package main

import (
	"bufio"
	"flag"
	"os"
)

func main() {
	pskip := flag.Int("n", 1, "skip first N lines")
	flag.Parse()

	skip := *pskip
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		if skip > 0 {
			skip--
			continue
		}

		os.Stdout.Write(s.Bytes())
		os.Stdout.Write([]byte{'\n'})
	}
}
