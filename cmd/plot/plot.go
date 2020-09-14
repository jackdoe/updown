package main

import (
	"bufio"
	"flag"
	"github.com/wcharczuk/go-chart" //exposes "chart"
	"github.com/wcharczuk/go-chart/drawing"
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
	var ptype = flag.String("mode", "line", "type line, scatter")
	var pmovingAverage = flag.Bool("line-moving-average", false, "also show moving average")
	flag.Parse()

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

	var graph chart.Chart
	if *ptype == "scatter" {
		viridisByY := func(xr, yr chart.Range, index int, x, y float64) drawing.Color {
			return chart.Viridis(y, yr.GetMin(), yr.GetMax())
		}

		graph = chart.Chart{
			XAxis: chart.XAxis{
				Style: chart.Shown(),
			},
			YAxis: chart.YAxis{
				Style:    chart.Shown(),
				AxisType: chart.YAxisSecondary,
			},
			Series: []chart.Series{
				chart.ContinuousSeries{
					Style: chart.Style{
						Hidden:           false,
						StrokeWidth:      chart.Disabled,
						DotWidth:         5,
						DotColorProvider: viridisByY,
					},
					XValues: x,
					YValues: y,
				},
			},
		}
	} else if *ptype == "line" {
		series := chart.ContinuousSeries{
			XValues: x,
			YValues: y,
		}

		graph = chart.Chart{
			XAxis: chart.XAxis{
				Style: chart.Shown(),
			},
			YAxis: chart.YAxis{
				Style:    chart.Shown(),
				AxisType: chart.YAxisSecondary,
			},
		}
		if *pmovingAverage {
			smaSeries := &chart.SMASeries{
				InnerSeries: series,
			}

			graph.Series = []chart.Series{
				series,
				smaSeries,
			}
		} else {
			graph.Series = []chart.Series{series}
		}
	} else {
		panic("unknown --mode")
	}
	err := graph.Render(chart.PNG, os.Stdout)
	if err != nil {
		panic(err)
	}
}
