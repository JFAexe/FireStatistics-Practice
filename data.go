package main

import (
	"os"
	"strconv"
	"strings"

	df "github.com/go-gota/gota/dataframe"
	sr "github.com/go-gota/gota/series"
	ch "github.com/vicanso/go-charts/v2"
)

func ReadDataFile(path string) df.DataFrame {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return df.ReadCSV(file, df.WithDelimiter(';'))
}

func ProcessData(path string) {
	basedf := ReadDataFile(path).Rename("id", "type_id").Rename("name", "type_name")

	dates := basedf.Select("dt").Records()[1:]

	basedf = basedf.Mutate(sr.New(Map(dates, func(i []string) int { return ParseDate(i).Year() }), sr.Int, "year"))
	basedf = basedf.Mutate(sr.New(Map(dates, func(i []string) int { return int(ParseDate(i).Month()) }), sr.Int, "month"))
	basedf = basedf.Mutate(sr.New(Map(dates, func(i []string) int { return ParseDate(i).Day() }), sr.Int, "day"))

	basedf = basedf.Drop("dt").Arrange(df.Sort("year"))

	years := basedf.Col("year")

	arr := []float64{}
	for i := years.Min(); i <= years.Max(); i++ {
		count := basedf.Filter(df.F{Colname: "year", Comparator: sr.Eq, Comparando: i}).Nrow()

		arr = append(arr, float64(count))
	}

	PrintMemoryUsage()

	cd := ChartData{
		title:  strings.Join([]string{"Количество пожаров за год. Общее число:", strconv.Itoa(basedf.Nrow())}, " "),
		values: [][]float64{arr},
		xaxis:  RemoveDuplicates(years.Records()),
		width:  640,
		height: 200,
		padding: ch.Box{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
	}

	if err := WriteFileFromBytes(path, ".svg", MakeBarChart(cd)); err != nil {
		panic(err)
	}
}

type ChartData struct {
	title   string
	values  [][]float64
	xaxis   []string
	width   int
	height  int
	padding ch.Box
}

func MakeBarChart(d ChartData) []byte {
	render, err := ch.BarRender(
		d.values,
		ch.XAxisDataOptionFunc(d.xaxis),
		ch.TitleTextOptionFunc(d.title),
		ch.WidthOptionFunc(d.width),
		ch.HeightOptionFunc(d.height),
		ch.PaddingOptionFunc(d.padding),
		ch.SVGTypeOption(),
	)
	if err != nil {
		panic(err)
	}

	data, err := render.Bytes()
	if err != nil {
		panic(err)
	}

	return data
}
