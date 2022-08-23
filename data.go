package main

import (
	"os"
	"sort"

	df "github.com/go-gota/gota/dataframe"
	sr "github.com/go-gota/gota/series"
)

var humanmonths = []string{
	"Январь",
	"Февраль",
	"Март",
	"Апрель",
	"Май",
	"Июнь",
	"Июль",
	"Август",
	"Сентябрь",
	"Октябрь",
	"Ноябрь",
	"Декабрь",
}

func ReadDataFile(path string) df.DataFrame {
	file, err := os.Open(path)
	if err != nil {
		ErrorLogger.Panic(err)
	}
	defer file.Close()

	return df.ReadCSV(file, df.WithDelimiter(';'))
}

func PrepareDataFrame(f df.DataFrame) df.DataFrame {
	dt := f.Select("dt").Records()[1:]

	return f.
		Mutate(sr.New(Map(dt, func(i []string) int { return ParseDate(i).Year() }), sr.Int, "year")).
		Mutate(sr.New(Map(dt, func(i []string) int { return int(ParseDate(i).Month()) }), sr.Int, "month")).
		Mutate(sr.New(Map(dt, func(i []string) int { return ParseDate(i).Day() }), sr.Int, "day")).
		Rename("type", "type_id").
		Rename("name", "type_name").
		Drop("dt").
		Arrange(df.Sort("year"))
}

func GetUniqueRecords(f df.DataFrame, n string) []string {
	return RemoveDuplicates(f.Col(n).Records())
}

func FilterEq(n, v string) df.F {
	return df.F{Colname: n, Comparator: sr.Eq, Comparando: v}
}

func CountFilteredData(f df.DataFrame, r []string, c string) []float64 {
	data := []float64{}
	for _, v := range r {
		data = append(data, float64(f.Filter(FilterEq(c, v)).Nrow()))
	}

	return data
}

func DoubleFilteredData(f df.DataFrame, r1, r2 []string, c1, c2 string) [][]float64 {
	data := [][]float64{}
	for _, v1 := range r1 {
		tf := f.Filter(FilterEq(c1, v1))
		td := []float64{}
		for _, v2 := range r2 {
			td = append(td, float64(tf.Filter(FilterEq(c2, v2)).Nrow()))
		}
		data = append(data, td)
	}

	return data
}

func GetMonthData(f df.DataFrame) []string {
	t := GetUniqueRecords(f, "month")

	sort.Slice(t, func(i, j int) bool { return ParseNumber(t[i]) < ParseNumber(t[j]) })

	return t
}

func ProcessData(path string) {
	frame := PrepareDataFrame(ReadDataFile(path))

	count := frame.Nrow()

	years := GetUniqueRecords(frame, "year")

	months := GetMonthData(frame)
	months_names := []string{}
	for _, m := range months {
		months_names = append(months_names, humanmonths[ParseNumber(m)-1])
	}

	types := GetUniqueRecords(frame, "type")
	types_names := []string{}
	for _, t := range types {
		types_names = append(types_names, frame.Filter(FilterEq("type", t)).Col("name").Records()[0])
	}

	InfoLogger.Println(years, count)
	InfoLogger.Println(months, months_names)
	InfoLogger.Println(types, types_names)

	title := []string{"Распределение за ", years[0], "-", years[len(years)-1]}

	count_years_total := CountFilteredData(frame, years, "year")
	count_total := DoubleFilteredData(frame, months, years, "month", "year")
	count_years := DoubleFilteredData(frame, years, months, "year", "month")

	GenerateBarChart(path, []string{"count_total_years.svg"}, [][]float64{count_years_total}, years, nil, []string{"Количество в год"}, minwidth)
	GenerateBarChart(path, []string{"count_total_span.svg"}, count_total, years, months_names, title, maxwidth)
	for i, y := range count_years {
		GenerateBarChart(path, []string{"count_", years[i], ".svg"}, [][]float64{y}, months_names, nil, []string{years[i]}, minwidth)
	}

	types_count_total := CountFilteredData(frame, types, "type")
	types_total := DoubleFilteredData(frame, types, years, "type", "year")
	types_years := DoubleFilteredData(frame, years, types, "year", "type")

	GeneratePieChart(path, []string{"types_total_percentage.svg"}, types_count_total, types_names, minwidth)
	GenerateBarChart(path, []string{"types_total_span.svg"}, types_total, years, types_names, title, defwidth)
	for i, y := range types_years {
		GeneratePieChart(path, []string{"types_", years[i], ".svg"}, y, types_names, minwidth)
	}

	GenerateBar(PrepareDataForCharts(count_years_total), years)
}
