package main

import (
	"os"
	"sort"
	"strconv"
	"strings"

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
	f, err := os.Open(path)
	if err != nil {
		ErrorLogger.Panic(err)
	}
	defer f.Close()

	return df.ReadCSV(f, df.WithDelimiter(';'))
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

func CountFilteredData(f df.DataFrame, r []string, c string) ([]int, [][][]float64) {
	d := make([]int, 0)
	p := make([][][]float64, 0)

	for _, v := range r {
		ff := f.Filter(FilterEq(c, v))

		d = append(d, ff.Nrow())
		p = append(p, Map(ff.Select([]string{"lon", "lat"}).Records()[1:], func(i []string) []float64 { return ParseFloatArray(i) }))
	}

	return d, p
}

func DoubleFilteredData(f df.DataFrame, r1, r2 []string, c1, c2 string) ([][]int, [][][][]float64) {
	d := make([][]int, 0)
	p := make([][][][]float64, 0)

	for _, v1 := range r1 {
		tf := f.Filter(FilterEq(c1, v1))

		td := make([]int, 0)
		tp := make([][][]float64, 0)

		for _, v2 := range r2 {
			tff := tf.Filter(FilterEq(c2, v2))

			td = append(td, tff.Nrow())
			tp = append(tp, Map(tff.Select([]string{"lon", "lat"}).Records()[1:], func(i []string) []float64 { return ParseFloatArray(i) }))
		}

		d = append(d, td)
		p = append(p, tp)
	}

	return d, p
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

	span := strings.Join([]string{years[0], "-", years[len(years)-1]}, "")

	months := GetMonthData(frame)
	months_names := make([]string, 0)
	for _, m := range months {
		months_names = append(months_names, humanmonths[ParseNumber(m)-1])
	}

	types := GetUniqueRecords(frame, "type")
	types_names := make([]string, 0)
	for _, t := range types {
		types_names = append(types_names, frame.Filter(FilterEq("type", t)).Col("name").Records()[0])
	}

	InfoLogger.Println(years, count)
	InfoLogger.Println(months, months_names)
	InfoLogger.Println(types, types_names)

	count_years_total, _ := CountFilteredData(frame, years, "year")
	count_total, _ := DoubleFilteredData(frame, months, years, "month", "year")
	count_years, _ := DoubleFilteredData(frame, years, months, "year", "month")

	types_count_total, _ := CountFilteredData(frame, types, "type")
	types_total, _ := DoubleFilteredData(frame, types, years, "type", "year")
	types_years, _ := DoubleFilteredData(frame, years, types, "year", "type")

	out := MakePage()

	out.AddCharts(
		geoBase([][][]float64{}),
		BarChart(strings.Join([]string{"Число за", span, "(", strconv.Itoa(count), ")"}, " "), years, count_years_total),
		PieChart(strings.Join([]string{"Отношение за", span}, " "), types_names, types_count_total),
		BarChartNestedValues(strings.Join([]string{"Распределение за", span}, " "), years, months_names, count_total),
		BarChartNestedValues(strings.Join([]string{"Распределение за", span}, " "), years, types_names, types_total),
	)
	for i := 0; i < len(count_years); i++ {
		out.AddCharts(
			BarChart(strings.Join([]string{"Число за", years[i]}, " "), months_names, count_years[i]),
			PieChart(strings.Join([]string{"Отношение за", years[i]}, " "), types_names, types_years[i]),
		)
	}

	RenderPage(out, GetFileNameFromPath(path), "out.html")
}
