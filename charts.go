package main

import (
	"strconv"
	"strings"

	ch "github.com/vicanso/go-charts/v2"
)

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

func MakeYearsTotalCountChart(path string, values []float64, axis []string, count int) {
	render, err := ch.BarRender(
		[][]float64{values},
		ch.TitleOptionFunc(ch.TitleOption{
			Text:            "Количество пожаров за год",
			Subtext:         strings.Join([]string{"Общее число:", strconv.Itoa(count)}, " "),
			Left:            ch.PositionCenter,
			FontSize:        24,
			SubtextFontSize: 16,
		}),
		ch.XAxisDataOptionFunc(axis),
		ch.WidthOptionFunc(60*len(axis)),
		ch.HeightOptionFunc(320),
		ch.PaddingOptionFunc(ch.Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}),
		ch.SVGTypeOption(),
	)
	if err != nil {
		ErrorLogger.Panic(err)
	}

	WriteChartFile(path, "count_years_total.svg", GetDataFromRender(render))
}

func MakeMonthsTotalCountChart(path string, values [][]float64, axis, legend []string) {
	render, err := ch.BarRender(
		values,
		ch.TitleOptionFunc(ch.TitleOption{
			Text:     "Месяц",
			Left:     ch.PositionCenter,
			FontSize: 24,
		}),
		ch.XAxisDataOptionFunc(axis),
		ch.LegendOptionFunc(ch.LegendOption{
			Data:   legend,
			Orient: ch.OrientHorizontal,
			Left:   ch.PositionCenter,
			Padding: ch.Box{
				Top:    48,
				Right:  0,
				Bottom: 0,
				Left:   0,
			},
		}),
		ch.WidthOptionFunc(140*len(axis)),
		ch.HeightOptionFunc(320),
		ch.PaddingOptionFunc(ch.Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}),
		ch.SVGTypeOption(),
	)
	if err != nil {
		ErrorLogger.Panic(err)
	}

	WriteChartFile(path, "count_months_total.svg", GetDataFromRender(render))
}

func MakeYearMonthsCountChart(path, year string, values [][]float64, axis, legend []string) {
	render, err := ch.BarRender(
		values,
		ch.TitleOptionFunc(ch.TitleOption{
			Text:     year,
			Left:     ch.PositionCenter,
			FontSize: 24,
		}),
		ch.XAxisDataOptionFunc(axis),
		ch.LegendOptionFunc(ch.LegendOption{
			Data:   legend,
			Orient: ch.OrientHorizontal,
			Left:   ch.PositionCenter,
			Padding: ch.Box{
				Top:    48,
				Right:  0,
				Bottom: 0,
				Left:   0,
			},
		}),
		ch.WidthOptionFunc(60*len(axis)),
		ch.HeightOptionFunc(320),
		ch.PaddingOptionFunc(ch.Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}),
		ch.SVGTypeOption(),
	)
	if err != nil {
		ErrorLogger.Panic(err)
	}

	WriteChartFile(path, strings.Join([]string{"count_", year, ".svg"}, ""), GetDataFromRender(render))
}

func MakeTypesTotalCountChart(path string, values [][]float64, axis, legend []string) {
	render, err := ch.BarRender(
		values,
		ch.TitleOptionFunc(ch.TitleOption{
			Text:     "Типы",
			Left:     ch.PositionCenter,
			FontSize: 24,
		}),
		ch.XAxisDataOptionFunc(axis),
		ch.LegendOptionFunc(ch.LegendOption{
			Data:   legend,
			Orient: ch.OrientHorizontal,
			Left:   ch.PositionCenter,
			Padding: ch.Box{
				Top:    48,
				Right:  0,
				Bottom: 0,
				Left:   0,
			},
		}),
		ch.WidthOptionFunc(110*len(axis)),
		ch.HeightOptionFunc(320),
		ch.PaddingOptionFunc(ch.Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}),
		ch.SVGTypeOption(),
	)
	if err != nil {
		ErrorLogger.Panic(err)
	}

	WriteChartFile(path, "types_total.svg", GetDataFromRender(render))
}
