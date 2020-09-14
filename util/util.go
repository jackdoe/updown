package util

import (
	"bufio"
	"io"
	"strconv"
)

// NL = new line
var NL = []byte{'\n'}

// ForeachLine also supporting lines not ending with \n a bit better than scanner.Scan()
func ForeachLine(in io.Reader, cb func(s string, b bool)) {
	r := bufio.NewReader(in)
	for {
		lineNL, err := r.ReadString('\n')
		if err == io.EOF {
			if len(lineNL) > 0 {
				cb(lineNL, false)
			}
			break
		}
		if err != nil {
			panic(err)
		}

		cb(lineNL[:len(lineNL)-1], true)
	}
}

// IntOrPanic die if not int
func IntOrPanic(s string) int64 {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// FloatOrPanic die if not float
func FloatOrPanic(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
