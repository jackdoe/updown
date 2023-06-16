package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

func main() {
	model := flag.String("m", "gpt-3.5-turbo", "options: gpt-4-32k, gpt-4, gpt-3.5-turbo, etc")
	print := flag.Bool("p", false, "print the tokens")
	flag.Parse()

	tkm, err := tiktoken.EncodingForModel(*model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	b, _ := ioutil.ReadAll(os.Stdin)
	tokens := tkm.Encode(string(b), nil, nil)
	if *print {
		for i, token := range tokens {
			os.Stdout.Write([]byte(fmt.Sprintf("%d", token)))
			if i == len(tokens)-1 {
				os.Stdout.Write([]byte("\n"))
			} else {
				os.Stdout.Write([]byte(" "))
			}
		}
	} else {
		fmt.Printf("%d\n", len(tokens))
	}
}
