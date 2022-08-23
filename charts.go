package main

import (
	"io"

	ch "github.com/go-echarts/go-echarts/v2/charts"
	cm "github.com/go-echarts/go-echarts/v2/components"
	op "github.com/go-echarts/go-echarts/v2/opts"
	tp "github.com/go-echarts/go-echarts/v2/types"
)

func MakePage(c ...cm.Charter) *cm.Page {
	page := cm.NewPage().SetLayout(cm.PageFlexLayout)

	page.AddCharts(c...)

	return page
}

func RenderPage(page *cm.Page, path, name string) {
	page.Render(io.MultiWriter(CreateFile(path, name)))
}

func WithRenderer(opt Renderer) ch.GlobalOpts {
	return func(bc *ch.BaseConfiguration) {
		bc.Renderer = opt
	}
}

func DefaultOptions(title string, chart interface{ Validate() }) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title, Left: "center"}),
		ch.WithInitializationOpts(op.Initialization{Width: "45vw", Height: "40vh"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
		WithRenderer(NewSnippetRenderer(chart, chart.Validate)),
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

func ConverDataGeo(names []string, data []ArrangedPoints) [][]op.GeoData {
	ret := make([][]op.GeoData, 0)

	for _, arrange := range data {
		for id, points := range arrange {
			ret2 := make([]op.GeoData, 0)
			for _, point := range FilterPoints(2, points) {
				ret2 = append(ret2, op.GeoData{Name: names[id], Value: []any{point[0], point[1], "count"}})
			}
			ret = append(ret, ret2)
		}
	}

	return ret
}

func BarChart(title string, axis []string, values []int) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	chart.AddSeries("", ConverDataBar(values))

	chart.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return chart
}

func BarChartNestedValues(title string, axis, names []string, values [][]int) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	for id, val := range values {
		chart.AddSeries(names[id], ConverDataBar(val))
	}

	chart.SetXAxis(axis)

	return chart
}

func PieChart(title string, axis []string, values []int) *ch.Pie {
	chart := ch.NewPie()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	chart.AddSeries("", ConverDataPie(axis, values))

	chart.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return chart
}

func GeoChart(names []string, data []ArrangedPoints) *ch.Geo {
	chart := ch.NewGeo()

	chart.SetGlobalOptions(
		ch.WithInitializationOpts(op.Initialization{Width: "90vw", Height: "70vh"}),
		ch.WithGeoComponentOpts(op.GeoComponent{Map: "Russia"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
		WithRenderer(NewSnippetRenderer(chart, chart.Validate)),
	)

	for _, series := range ConverDataGeo(names, data) {
		chart.AddSeries("", tp.ChartScatter, series)
	}

	return chart
}
