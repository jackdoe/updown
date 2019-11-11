package main

import (
	"bufio"
	"flag"
	"fmt" //exposes "chart"
	"os"
	"strconv"
	"strings"

	"github.com/bmizerany/perks/quantile"
)

func floatOrPanic(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	var pop = flag.String("filter-op", "", "filter opt: gt, lt")
	var fvalue = flag.Float64("filter-value", 0, "use values lt or gt (depending on filter-op)")
	var fquant = flag.String("quantiles", "0.25,0.50,0.75,0.90,0.99", "which quantiles")
	flag.Parse()
	filterValue := *fvalue
	filter := func(v float64) bool {
		return true
	}
	if *pop == "gt" {
		filter = func(v float64) bool {
			return v > filterValue
		}
	}
	if *pop == "lt" {
		filter = func(v float64) bool {
			return v < filterValue
		}
	}

	s := bufio.NewScanner(os.Stdin)
	quantiles := []float64{}
	for _, v := range strings.Split(*fquant, ",") {
		quantiles = append(quantiles, floatOrPanic(v))
	}
	q := quantile.NewTargeted(quantiles...)
	for s.Scan() {
		text := s.Text()
		v := floatOrPanic(text)
		if filter(v) {
			q.Insert(v)
		}
	}
	for _, quant := range quantiles {
		fmt.Printf("perc%d: %f\n", int(quant*100), q.Query(quant))
	}

	fmt.Println("count:", q.Count())
}
