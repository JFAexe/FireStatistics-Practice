package main

import (
	"io"

	om "github.com/elliotchance/orderedmap/v2"
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
	return func(bc *ch.BaseConfiguration) { bc.Renderer = opt }
}

func WithScatterSize(opt float32) ch.SeriesOpts {
	return func(s *ch.SingleSeries) { s.SymbolSize = opt }
}

func DefaultOptions(title string, chart interface{ Validate() }) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title, Left: "center"}),
		ch.WithInitializationOpts(op.Initialization{Width: "45vw", Height: "40vh"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
		WithRenderer(NewSnippetRenderer(chart, chart.Validate)),
	}
}

func ConverDataBar(data om.OrderedMap[string, int]) ([]op.BarData, []string) {
	ret := make([]op.BarData, data.Len())
	axs := make([]string, data.Len())

	for id, key := range data.Keys() {
		value, _ := data.Get(key)
		ret[id] = op.BarData{Name: key, Value: value}
		axs[id] = key
	}

	return ret, axs
}

func BarChart(title string, values om.OrderedMap[string, int]) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	series, axis := ConverDataBar(values)

	chart.AddSeries("", series)

	chart.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return chart
}

func BarChartNestedValues(title string, axis []string, values om.OrderedMap[string, om.OrderedMap[string, int]]) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	for _, key := range values.Keys() {
		data, _ := values.Get(key)
		series, _ := ConverDataBar(data)
		chart.AddSeries(key, series)
	}

	chart.SetXAxis(axis)

	return chart
}

func ConverDataPie(data om.OrderedMap[string, int]) []op.PieData {
	ret := make([]op.PieData, data.Len())

	for id, key := range data.Keys() {
		value, _ := data.Get(key)
		ret[id] = op.PieData{Name: key, Value: value}
	}

	return ret
}

func PieChart(title string, values om.OrderedMap[string, int]) *ch.Pie {
	chart := ch.NewPie()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	chart.AddSeries("", ConverDataPie(values))

	chart.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return chart
}

func ConverDataGeo(data Points, tip string) ([]op.GeoData, float32, float32) {
	ret := make([]op.GeoData, 0)

	min, max := 1, 1

	for point, count := range FilterPoints(rad, data) {
		if count < min {
			min = count
		}

		if count > max {
			max = count
		}

		ret = append(ret, op.GeoData{Name: tip, Value: []any{point.x, point.y, count}})
	}

	return ret, float32(min), float32(max)
}

func GeoChart(title string, data om.OrderedMap[string, Points]) *ch.Geo {
	chart := ch.NewGeo()

	var cmin, cmax float32 = 1, 1

	for _, key := range data.Keys() {
		points, _ := data.Get(key)
		series, min, max := ConverDataGeo(points, key)

		if cmin < cmin {
			cmin = min
		}

		if max > cmax {
			cmax = max
		}

		chart.AddSeries("", tp.ChartScatter, series, WithScatterSize(dia))
	}

	chart.SetGlobalOptions(append(
		DefaultOptions(title, chart),
		ch.WithGeoComponentOpts(op.GeoComponent{Map: "Russia"}),
		ch.WithInitializationOpts(op.Initialization{Width: "90vw", Height: "70vh"}),
	)...)

	return chart
}
