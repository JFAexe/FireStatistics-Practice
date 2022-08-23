package main

import (
	"io"

	ch "github.com/go-echarts/go-echarts/v2/charts"
	cm "github.com/go-echarts/go-echarts/v2/components"
	op "github.com/go-echarts/go-echarts/v2/opts"
	rn "github.com/go-echarts/go-echarts/v2/render"
	tp "github.com/go-echarts/go-echarts/v2/types"
)

func WithRenderer(opt rn.Renderer) ch.GlobalOpts {
	return func(bc *ch.BaseConfiguration) {
		bc.Renderer = opt
	}
}

func MakePage() *cm.Page {
	return cm.NewPage().SetLayout(cm.PageFlexLayout)
}

func RenderPage(page *cm.Page, path, name string) {
	page.Render(io.MultiWriter(CreateFile(path, name)))
}

func DefaultOptions(title string, renderer rn.Renderer) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title, Left: "center"}),
		ch.WithInitializationOpts(op.Initialization{Width: "45vw", Height: "40vh"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
		WithRenderer(renderer),
	}
}

func ConverDataBar(data []int) []op.BarData {
	ret := make([]op.BarData, len(data))

	for id, val := range data {
		ret[id] = op.BarData{Value: val}
	}

	return ret
}

func ConverDataPie(names []string, data []int) []op.PieData {
	ret := make([]op.PieData, len(data))

	for id, val := range data {
		ret[id] = op.PieData{Name: names[id], Value: val}
	}

	return ret
}

func ConverDataGeo(data []ArrangedPoints) []op.GeoData {
	ret := make([]op.GeoData, 0)

	count := 0

	for _, arrange := range data {
		for _, points := range arrange {
			for _, point := range points {
				if count > 1023 {
					break
				}
				ret = append(ret, op.GeoData{Value: point})
				count++
			}
		}
	}

	return ret
}

func BarChart(title string, axis []string, values []int) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, NewSnippetRenderer(chart, chart.Validate))...)

	chart.AddSeries("", ConverDataBar(values))

	chart.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return chart
}

func BarChartNestedValues(title string, axis, names []string, values [][]int) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, NewSnippetRenderer(chart, chart.Validate))...)

	for id, val := range values {
		chart.AddSeries(names[id], ConverDataBar(val))
	}

	chart.SetXAxis(axis)

	return chart
}

func PieChart(title string, axis []string, values []int) *ch.Pie {
	chart := ch.NewPie()

	chart.SetGlobalOptions(DefaultOptions(title, NewSnippetRenderer(chart, chart.Validate))...)

	chart.AddSeries("", ConverDataPie(axis, values))

	chart.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return chart
}

func GeoChart(data []ArrangedPoints) *ch.Geo {
	chart := ch.NewGeo()

	chart.SetGlobalOptions(
		ch.WithInitializationOpts(op.Initialization{Width: "90vw", Height: "70vh"}),
		ch.WithGeoComponentOpts(op.GeoComponent{Map: "Russia"}),
		WithRenderer(NewSnippetRenderer(chart, chart.Validate)),
	)

	chart.AddSeries("", tp.ChartEffectScatter, ConverDataGeo(data))

	return chart
}
