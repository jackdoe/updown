package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Aggregate struct {
	StringValue string
	FloatValue  float64
	Position    int
}

func main() {
	column := flag.Int("col-index", 0, "column to group by")
	sep := flag.String("separator", ",", "separator")
	flag.Parse()

	r := csv.NewReader(os.Stdin)
	separator := *sep
	r.Comma = rune(separator[0])

	byDate := map[string][]*Aggregate{}
	position := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		key := record[*column]
		agg, ok := byDate[key]
		if !ok {
			agg = make([]*Aggregate, len(record))
			byDate[key] = agg
		}

		for idx, v := range record {
			aggregate := agg[idx]

			if aggregate == nil {
				aggregate = &Aggregate{Position: position, StringValue: v}
				agg[idx] = aggregate
			}
			if idx != *column {
				f, err := strconv.ParseFloat(v, 64)
				if err == nil {
					aggregate.FloatValue += f
					aggregate.StringValue = ""
				}
			}
		}
	}

	for _, v := range byDate {
		for i, a := range v {
			if a.StringValue != "" {
				os.Stdout.Write([]byte(a.StringValue))
			} else {
				os.Stdout.Write([]byte(fmt.Sprintf("%.2f", a.FloatValue)))
			}
			if i != len(v)-1 {
				os.Stdout.Write([]byte(separator))
			}
		}
		os.Stdout.Write([]byte{10})
	}
}
