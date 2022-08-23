package main

import (
	"io"

	ch "github.com/go-echarts/go-echarts/v2/charts"
	cm "github.com/go-echarts/go-echarts/v2/components"
	op "github.com/go-echarts/go-echarts/v2/opts"
	tp "github.com/go-echarts/go-echarts/v2/types"
)

func MakePage() *cm.Page {
	return cm.NewPage().SetLayout(cm.PageFlexLayout)
}

func RenderPage(page *cm.Page, path, name string) {
	page.Render(io.MultiWriter(CreateFile(path, name)))
}

func DefaultOptions(title string) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title, Left: "center"}),
		ch.WithInitializationOpts(op.Initialization{Width: "45vw", Height: "40vh"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
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
	ret := make([]op.GeoData, 1024)

	count := 0

	for _, arrange := range data {
		for _, points := range arrange {
			for _, point := range points {
				if count > 1023 {
					break
				}
				ret[count] = op.GeoData{Value: point}
				count++
			}
		}
	}

	return ret
}

func BarChart(title string, axis []string, values []int) *ch.Bar {
	bar := ch.NewBar()

	bar.SetGlobalOptions(DefaultOptions(title)...)

	bar.AddSeries("", ConverDataBar(values))

	bar.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return bar
}

func BarChartNestedValues(title string, axis, names []string, values [][]int) *ch.Bar {
	bar := ch.NewBar()

	bar.SetGlobalOptions(DefaultOptions(title)...)

	for id, val := range values {
		bar.AddSeries(names[id], ConverDataBar(val))
	}

	bar.SetXAxis(axis)

	return bar
}

func PieChart(title string, axis []string, values []int) *ch.Pie {
	pie := ch.NewPie()

	pie.SetGlobalOptions(DefaultOptions(title)...)

	pie.AddSeries("", ConverDataPie(axis, values))

	pie.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return pie
}

func GeoChart(data []ArrangedPoints) *ch.Geo {
	geo := ch.NewGeo()

	geo.SetGlobalOptions(
		ch.WithInitializationOpts(op.Initialization{Width: "90vw", Height: "70vh"}),
		ch.WithGeoComponentOpts(op.GeoComponent{Map: "Russia"}),
	)

	geo.AddSeries("", tp.ChartEffectScatter, ConverDataGeo(data))

	return geo
}
