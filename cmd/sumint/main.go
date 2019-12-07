package main

import (
	"bufio"
	"fmt" //exposes "chart"
	"os"
	"strconv"
)

func intOrPanic(s string) int64 {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := int64(0)
	for s.Scan() {
		text := s.Text()
		v := intOrPanic(text)
		sum += v
	}
	fmt.Printf("%d\n", sum)
}
