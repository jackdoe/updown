package main

import (
	"bufio"
	"github.com/wcharczuk/go-chart" //exposes "chart"
	"os"
	"regexp"
	"strconv"
)

func floatOrPanic(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	i := 0
	x := []float64{}
	y := []float64{}

	re := regexp.MustCompile("\\s")

	for s.Scan() {
		text := s.Text()
		splitted := re.Split(text, -1)

		if len(splitted) == 2 {
			x = append(x, floatOrPanic(splitted[0]))
			y = append(y, floatOrPanic(splitted[1]))
		} else {
			y = append(y, floatOrPanic(splitted[0]))
			x = append(x, float64(i))
		}
		i++
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: y,
			},
		},
	}

	err := graph.Render(chart.PNG, os.Stdout)
	if err != nil {
		panic(err)
	}
}
