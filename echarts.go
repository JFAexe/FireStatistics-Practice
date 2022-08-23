package main

import (
	"io"

	ch "github.com/go-echarts/go-echarts/v2/charts"
	cm "github.com/go-echarts/go-echarts/v2/components"
	op "github.com/go-echarts/go-echarts/v2/opts"
)

func MakePage() *cm.Page {
	return cm.NewPage()
}

func RenderPage(p *cm.Page, path, name string) {
	p.Render(io.MultiWriter(CreateFile(path, name)))
}

func ConverDataBar[T any](d []T) []op.BarData {
	items := make([]op.BarData, 0)

	for _, v := range d {
		items = append(items, op.BarData{Value: v})
	}

	return items
}

func ConverDataPie[T any](n []string, d []T) []op.PieData {
	items := make([]op.PieData, 0)

	for i, v := range d {
		items = append(items, op.PieData{Name: n[i], Value: v})
	}

	return items
}

func BarChart[T any](title string, axis []string, values ...[]T) *ch.Bar {
	bar := ch.NewBar()

	bar.SetGlobalOptions(DefaultOptions(title)...)

	for i := 0; i < len(values); i++ {
		bar.AddSeries("", ConverDataBar(values[i]))
	}

	bar.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return bar
}

func BarChartNestedValues[T any](title string, axis, names []string, values ...[][]T) *ch.Bar {
	bar := ch.NewBar()

	bar.SetGlobalOptions(DefaultOptions(title)...)

	for i := 0; i < len(values); i++ {
		for j, v := range values[i] {
			bar.AddSeries(names[j], ConverDataBar(v))
		}
	}

	bar.SetXAxis(axis)

	return bar
}

func PieChart[T any](title string, axis []string, values []T) *ch.Pie {
	pie := ch.NewPie()

	pie.SetGlobalOptions(DefaultOptions(title)...)

	pie.AddSeries("", ConverDataPie(axis, values))

	pie.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return pie
}

func DefaultOptions(title string) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title}),
		ch.WithInitializationOpts(op.Initialization{Width: "960px", Height: "320px"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
	}
}
