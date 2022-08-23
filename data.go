package main

import (
	"os"
	"strconv"
	"strings"

	df "github.com/go-gota/gota/dataframe"
	sr "github.com/go-gota/gota/series"
	ch "github.com/vicanso/go-charts/v2"
)

func FilterEq(n, v string) df.F {
	return df.F{Colname: n, Comparator: sr.Eq, Comparando: v}
}

func ReadDataFile(path string) df.DataFrame {
	file, err := os.Open(path)
	if err != nil {
		ErrorLogger.Panic(err)
	}
	defer file.Close()

	return df.ReadCSV(file, df.WithDelimiter(';'))
}

func ProcessData(path string) {
	basedf := ReadDataFile(path)

	dates := basedf.Select("dt").Records()[1:]

	basedf = basedf.
		Mutate(sr.New(Map(dates, func(i []string) int { return ParseDate(i).Year() }), sr.Int, "year")).
		Mutate(sr.New(Map(dates, func(i []string) int { return int(ParseDate(i).Month()) }), sr.Int, "month")).
		Mutate(sr.New(Map(dates, func(i []string) int { return ParseDate(i).Day() }), sr.Int, "day"))

	basedf = basedf.
		Rename("type", "type_id").
		Rename("name", "type_name").
		Drop("dt").
		Arrange(df.Sort("year"))

	years := RemoveDuplicates(basedf.Col("year").Records()[1:])
	types := RemoveDuplicates(basedf.Col("type").Records()[1:])
	names := []string{}
	for _, t := range types {
		names = append(names, basedf.Filter(FilterEq("type", t)).Col("name").Records()[0])
	}

	years_counts := []float64{}
	for _, y := range years {
		count := basedf.Filter(FilterEq("year", y)).Nrow()

		years_counts = append(years_counts, float64(count))
	}

	years_types := [][]float64{}
	for _, t := range types {
		yt := []float64{}

		for _, y := range years {
			d := float64(basedf.Filter(FilterEq("year", y), FilterEq("type", t)).Nrow())
			yt = append(yt, d)
		}

		years_types = append(years_types, yt)
	}

	ytch := ChartData{
		title:  "Типы",
		values: years_types,
		axis:   years,
		legend: names,
		width:  100 * len(years),
		height: 350,
		padding: ch.Box{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
	}

	if err := WriteFileFromBytes(path, "types.svg", MakeVerticalBarChart(ytch)); err != nil {
		ErrorLogger.Panic(err)
	}

	total := ChartData{
		title:  strings.Join([]string{"Количество пожаров за год. Общее число:", strconv.Itoa(basedf.Nrow())}, " "),
		values: [][]float64{years_counts},
		axis:   years,
		legend: nil,
		width:  58 * len(years),
		height: 350,
		padding: ch.Box{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
	}

	if err := WriteFileFromBytes(path, "total.svg", MakeVerticalBarChart(total)); err != nil {
		ErrorLogger.Panic(err)
	}
}

type ChartData struct {
	values  [][]float64
	title   string
	axis    []string
	legend  []string
	width   int
	height  int
	padding ch.Box
}

func MakeVerticalBarChart(d ChartData) []byte {
	render, err := ch.BarRender(
		d.values,
		ch.TitleTextOptionFunc(d.title),
		ch.XAxisDataOptionFunc(d.axis),
		ch.LegendLabelsOptionFunc(d.legend),
		ch.WidthOptionFunc(d.width),
		ch.HeightOptionFunc(d.height),
		ch.PaddingOptionFunc(d.padding),
		ch.SVGTypeOption(),
	)
	if err != nil {
		ErrorLogger.Panic(err)
	}

	data, err := render.Bytes()
	if err != nil {
		ErrorLogger.Panic(err)
	}

	return data
}
