package main

import (
	"os"
	"sort"
	"strconv"
	"strings"

	df "github.com/go-gota/gota/dataframe"
	sr "github.com/go-gota/gota/series"
)

type (
	Point     struct{ x, y float64 }
	Points    []Point
	PointsMap map[string]Points
)

var (
	humanmonths = map[string]string{
		"1":  "Январь",
		"2":  "Февраль",
		"3":  "Март",
		"4":  "Апрель",
		"5":  "Май",
		"6":  "Июнь",
		"7":  "Июль",
		"8":  "Август",
		"9":  "Сентябрь",
		"10": "Октябрь",
		"11": "Ноябрь",
		"12": "Декабрь",
	}
)

func ReadDataFile(path string) df.DataFrame {
	file, err := os.Open(path)
	if err != nil {
		ErrorLogger.Panic(err)
	}
	defer file.Close()

	return df.ReadCSV(file, df.WithDelimiter(';'))
}

func PrepareDataFrame(frame df.DataFrame) df.DataFrame {
	dt := frame.Select("dt").Records()[1:]

	return frame.
		Rename("type", "type_id").
		Rename("name", "type_name").
		Mutate(sr.New(Map(dt, DateYear), sr.Int, "year")).
		Mutate(sr.New(Map(dt, DateMonth), sr.Int, "month")).
		Mutate(sr.New(Map(dt, DateDay), sr.Int, "day")).
		Drop("dt").
		Arrange(df.Sort("year"))
}

func FilterEq(name, value string) df.F {
	return df.F{Colname: name, Comparator: sr.Eq, Comparando: value}
}

func GetUniqueInts(frame df.DataFrame, name string) []string {
	ret, err := frame.Col(name).Int()
	if err != nil {
		ErrorLogger.Panicf("Can't parse int. Error: %s\n", err)
	}

	ret = RemoveDuplicateValues(ret)

	sort.Ints(ret)

	return Map(ret, func(in int) string { return strconv.Itoa(in) })
}

func GetPoints(frame df.DataFrame) Points {
	ret := make(Points, 0)

	lon := frame.Col("lon").Float()
	lat := frame.Col("lat").Float()

	for id := 0; id < len(lon); id++ {
		ret = append(ret, Point{lon[id], lat[id]})
	}

	return ret
}

func SingleFilterPass(frame df.DataFrame, r []string, c string) ([]int, PointsMap) {
	data := make([]int, 0)
	points := make(PointsMap, 0)

	for id, v := range r {
		f := frame.Filter(FilterEq(c, v))

		data = append(data, f.Nrow())
		points[r[id]] = GetPoints(f)
	}

	return data, points
}

func DoubleFilterPass(frame df.DataFrame, r1, r2 []string, c1, c2 string) ([][]int, []PointsMap) {
	data := make([][]int, 0)
	points := make([]PointsMap, 0)

	for _, v1 := range r1 {
		d, p := SingleFilterPass(frame.Filter(FilterEq(c1, v1)), r2, c2)

		data = append(data, d)
		points = append(points, p)
	}

	return data, points
}

func ProcessData(path string) {
	frame := PrepareDataFrame(ReadDataFile(path))

	count := frame.Nrow()

	years := GetUniqueInts(frame, "year")

	span := strings.Join([]string{years[0], "-", years[len(years)-1]}, "")

	months := GetUniqueInts(frame, "month")
	months_names := make([]string, 0)
	for _, m := range months {
		months_names = append(months_names, humanmonths[m])
	}

	types := GetUniqueInts(frame, "type")
	types_names := make([]string, 0)
	for _, t := range types {
		types_names = append(types_names, frame.Filter(FilterEq("type", t)).Col("name").Records()[0])
	}

	InfoLogger.Println(years, "Всего:", count)
	InfoLogger.Println(months, months_names)
	InfoLogger.Println(types, types_names)

	count_years_total, points1 := SingleFilterPass(frame, years, "year")
	count_total, _ := DoubleFilterPass(frame, months, years, "month", "year")
	count_years, _ := DoubleFilterPass(frame, years, months, "year", "month")

	types_count_total, _ := SingleFilterPass(frame, types, "type")
	types_total, _ := DoubleFilterPass(frame, types, years, "type", "year")
	types_years, _ := DoubleFilterPass(frame, years, types, "year", "type")

	out := MakePage(
		GeoChart([]PointsMap{points1}),
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

	RenderPage(out, GetFileNameFromPath(path), "page.html")
}
