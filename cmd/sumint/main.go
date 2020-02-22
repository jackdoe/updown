package main

import (
	"bufio"
	"flag"
	"fmt"
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

func floatOrPanic(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	isFloat := flag.Bool("f", false, "expect floating point numbers")
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)
	if !*isFloat {
		sum := int64(0)
		for s.Scan() {
			text := s.Text()
			v := intOrPanic(text)
			sum += v
		}
		fmt.Printf("%d\n", sum)
	} else {
		sum := float64(0)
		for s.Scan() {
			text := s.Text()
			v := floatOrPanic(text)
			sum += v
		}
		fmt.Printf("%f\n", sum)

	}
}
