package main

import (
	"encoding/json"
	"os"
)

func main() {
	jd := json.NewDecoder(os.Stdin)
	je := json.NewEncoder(os.Stdout)
	var v interface{}
	err := jd.Decode(&v)
	if err != nil {
		panic(err)
	}
	je.SetIndent("", "  ")
	err = je.Encode(v)
	if err != nil {
		panic(err)
	}
}
