package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dcadenas/pagerank"
)

func main() {
	follow := flag.Float64("prob-follow", 0.85, "the bigger the number, less probability we have to teleport to some random link")
	tolerance := flag.Float64("tolerance", 0.0001, "the smaller the number, the more exact the result will be but more CPU cycles will be needed.")
	flag.Parse()

	sid := map[string]int{}
	back := map[int]string{}

	graph := pagerank.New()

	s := bufio.NewScanner(os.Stdin)

	find := func(s string) int {
		v, ok := sid[s]
		if ok {
			return v
		}
		id := len(sid) + 1
		sid[s] = id
		back[id] = s
		return id
	}

	for s.Scan() {
		splitted := filter(strings.Split(s.Text(), " "))
		if len(splitted) < 2 {
			continue
		}

		from := find(splitted[0])
		for i := 1; i < len(splitted); i++ {
			to := find(splitted[i])
			graph.Link(from, to)
		}
	}

	graph.Rank(*follow, *tolerance, func(identifier int, rank float64) {
		fmt.Printf("%.4f %s\n", rank, back[identifier])
	})
}

func filter(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return out
}
