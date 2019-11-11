package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"os"
)

func main() {
	perLine := flag.Bool("l", false, "is there one json document one per line?")
	flag.Parse()

	je := json.NewEncoder(os.Stdout)
	je.SetIndent("", "  ")
	if !*perLine {
		jd := json.NewDecoder(os.Stdin)
		var v interface{}
		err := jd.Decode(&v)
		if err != nil {
			panic(err)
		}
		err = je.Encode(v)
		if err != nil {
			panic(err)
		}

	} else {
		r := bufio.NewReader(os.Stdin)
		for {
			data, err := r.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			var v interface{}
			err = json.Unmarshal(data, &v)
			if err != nil {
				panic(err)
			}
			err = je.Encode(v)
			if err != nil {
				panic(err)
			}
		}
	}

}
