package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	yd := yaml.NewDecoder(os.Stdin)
	ye := yaml.NewEncoder(os.Stdout)
	var v interface{}
	err := yd.Decode(&v)
	if err != nil {
		panic(err)
	}
	err = ye.Encode(v)
	if err != nil {
		panic(err)
	}
}
