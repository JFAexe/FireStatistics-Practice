package main

import (
	"io"
	"math/rand"

	ch "github.com/go-echarts/go-echarts/v2/charts"
	cm "github.com/go-echarts/go-echarts/v2/components"
	op "github.com/go-echarts/go-echarts/v2/opts"
	tp "github.com/go-echarts/go-echarts/v2/types"
)

func MakePage() *cm.Page {
	return cm.NewPage().SetLayout(cm.PageFlexLayout)
}

func RenderPage(p *cm.Page, path, name string) {
	p.Render(io.MultiWriter(CreateFile(path, name)))
}

func ConverDataBar(d []int) []op.BarData {
	items := make([]op.BarData, 0)

	for _, v := range d {
		items = append(items, op.BarData{Value: v})
	}

	return items
}

func ConverDataPie(n []string, d []int) []op.PieData {
	items := make([]op.PieData, 0)

	for i, v := range d {
		items = append(items, op.PieData{Name: n[i], Value: v})
	}

	return items
}

func BarChart(title string, axis []string, values ...[]int) *ch.Bar {
	bar := ch.NewBar()

	bar.SetGlobalOptions(DefaultOptions(title)...)

	for i := 0; i < len(values); i++ {
		bar.AddSeries("", ConverDataBar(values[i]))
	}

	bar.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return bar
}

func BarChartNestedValues(title string, axis, names []string, values ...[][]int) *ch.Bar {
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

func PieChart(title string, axis []string, values []int) *ch.Pie {
	pie := ch.NewPie()

	pie.SetGlobalOptions(DefaultOptions(title)...)

	pie.AddSeries("", ConverDataPie(axis, values))

	pie.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return pie
}

func DefaultOptions(title string) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title, Left: "center"}),
		ch.WithInitializationOpts(op.Initialization{Width: "45vw", Height: "40vh"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
	}
}

func geoData(d [][][]float64) []op.GeoData {
	points := make([]op.GeoData, 0)

	c := 0

	for _, v := range d {

		for _, p := range v {
			if c > 50 {
				break
			}
			points = append(points, op.GeoData{Value: []float64{p[0], p[1], float64(rand.Intn(100))}})
			c++
		}
	}

	return points
}

func geoBase(p [][][]float64) *ch.Geo {
	geo := ch.NewGeo()
	geo.SetGlobalOptions(
		ch.WithInitializationOpts(op.Initialization{Width: "90vw", Height: "70vh"}),
		ch.WithGeoComponentOpts(op.GeoComponent{Map: "Russia"}),
	)

	geo.AddSeries("geo", tp.ChartEffectScatter, geoData(p))

	return geo
}
