package main

import (
	"strings"

	ch "github.com/vicanso/go-charts/v2"
)

const (
	maxwidth  int = 1280
	defwidth  int = 940
	minwidth  int = 640
	defheight int = 360
)

var (
	basepadding   ch.Box = ch.Box{Top: 32, Right: 32, Bottom: 32, Left: 32}
	legendpadding ch.Box = ch.Box{Top: 56, Right: 0, Bottom: 0, Left: 0}
)

func TitleOptions(t string) ch.TitleOption {
	return ch.TitleOption{Text: t, Left: ch.PositionCenter, FontSize: 24}
}

func LegendOptions(d []string) ch.LegendOption {
	return ch.LegendOption{Data: d, Left: ch.PositionCenter, Padding: legendpadding, FontSize: 10}
}

func PieRadius(r string) ch.OptionFunc { // rekt
	return func(opt *ch.ChartOption) {
		for index := range opt.SeriesList {
			opt.SeriesList[index].Radius = r
		}
	}
}

func GetDataFromRender(c *ch.Painter) []byte {
	data, err := c.Bytes()
	if err != nil {
		ErrorLogger.Panic(err)
	}

	return data
}

func WriteChartFile(path, name string, data []byte) {
	dir := strings.Join([]string{temppath, GetFileNameFromPath(path)}, "/")

	if err := WriteFileFromBytes(dir, name, data); err != nil {
		ErrorLogger.Panic(err)
	}
}

func GenerateBarChart(path string, name []string, values [][]float64, axis, legend, title []string, width int) {
	render, err := ch.BarRender(
		values,
		ch.TitleOptionFunc(TitleOptions(strings.Join(title, ""))),
		ch.XAxisDataOptionFunc(axis),
		ch.LegendOptionFunc(LegendOptions(legend)),
		ch.WidthOptionFunc(width),
		ch.HeightOptionFunc(defheight),
		ch.PaddingOptionFunc(basepadding),
		ch.SVGTypeOption(),
	)
	if err != nil {
		ErrorLogger.Panic(err)
	}

	WriteChartFile(path, strings.Join(name, ""), GetDataFromRender(render))
}

func GeneratePieChart(path string, name []string, values []float64, legend []string, width int) {
	render, err := ch.PieRender(
		values,
		ch.LegendOptionFunc(ch.LegendOption{
			Data: legend,
			Show: ch.FalseFlag(),
		}),
		PieRadius("35%"),
		ch.PieSeriesShowLabel(),
		ch.WidthOptionFunc(width),
		ch.HeightOptionFunc(defheight),
		ch.PaddingOptionFunc(basepadding),
		ch.SVGTypeOption(),
	)
	if err != nil {
		ErrorLogger.Panic(err)
	}

	WriteChartFile(path, strings.Join(name, ""), GetDataFromRender(render))
}
