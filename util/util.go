package util

import (
	"bufio"
	"io"
)

func ForeachLine(in io.Reader, cb func(s string, b bool)) {
	r := bufio.NewReader(in)
	for {
		lineNL, err := r.ReadString('\n')
		if err == io.EOF {
			if len(lineNL) > 0 {
				cb(lineNL, true)
			}
			break
		}
		if err != nil {
			panic(err)
		}

		cb(lineNL, false)
	}
}
