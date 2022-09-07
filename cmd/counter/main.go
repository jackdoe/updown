package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	fn := path.Join(home, ".updown_counter.txt")
	b, err := ioutil.ReadFile(fn)
	v := uint64(0)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		v, _ = strconv.ParseUint(string(b), 10, 64)
	}
	fmt.Printf("%d\n", v)
	v++
	err = ioutil.WriteFile(fn, []byte(fmt.Sprintf("%d", v)), 0600)
	if err != nil {
		panic(err)
	}
}
