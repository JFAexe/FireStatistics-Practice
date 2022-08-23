package main

import (
	"io"
	"os"

	charts "github.com/go-echarts/go-echarts/v2/charts"
	components "github.com/go-echarts/go-echarts/v2/components"
	opts "github.com/go-echarts/go-echarts/v2/opts"
)

func PrepareDataForCharts[T any](d []T) []opts.BarData {
	items := make([]opts.BarData, 0)

	for _, v := range d {
		items = append(items, opts.BarData{Value: v})
	}

	return items
}

func barBasic(values []opts.BarData, axis []string) *charts.Bar {
	bar := charts.NewBar()

	bar.SetXAxis(axis).AddSeries("", values)

	return bar
}

func GenerateBar(values []opts.BarData, axis []string) {
	page := components.NewPage()
	page.AddCharts(barBasic(values, axis))
	f, err := os.Create("./temp/bar.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
