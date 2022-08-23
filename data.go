package main

import (
	"html/template"
	"os"
	"sort"
	"strings"

	om "github.com/elliotchance/orderedmap/v2"
	df "github.com/go-gota/gota/dataframe"
	sr "github.com/go-gota/gota/series"
)

type (
	Point  struct{ x, y float64 }
	Points []Point
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
		ErrorLogger.Fatal(err)
	}
	defer file.Close()

	return df.ReadCSV(file, df.WithDelimiter(';'))
}

func PrepareDataFrame(frame df.DataFrame) df.DataFrame {
	dt := frame.Select("dt").Records()[1:]

	return frame.
		Drop("dt").
		Rename("type", "type_id").
		Rename("name", "type_name").
		Mutate(sr.New(Map(dt, DateYear), sr.Int, "year")).
		Mutate(sr.New(Map(dt, DateMonth), sr.Int, "month")).
		Arrange(df.Sort("year"))
}

func FilterEq(name, value string) df.F {
	return df.F{Colname: name, Comparator: sr.Eq, Comparando: value}
}

func GetUniqueInts(frame df.DataFrame, name string) []string {
	ret, err := frame.Col(name).Int()
	if err != nil {
		ErrorLogger.Fatalf("Can't parse int. Error: %s\n", err)
	}

	ret = RemoveDuplicateValues(ret)

	sort.Ints(ret)

	return Map(ret, IntToStr)
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

func SingleFilterPass(frame df.DataFrame, r []string, c string) (om.OrderedMap[string, int], om.OrderedMap[string, Points]) {
	data := om.NewOrderedMap[string, int]()
	points := om.NewOrderedMap[string, Points]()

	for id, v := range r {
		f := frame.Filter(FilterEq(c, v))

		data.Set(r[id], f.Nrow())
		points.Set(r[id], GetPoints(f))
	}

	return *data, *points
}

func DoubleFilterPass(frame df.DataFrame, r1, r2 []string, c1, c2 string) (om.OrderedMap[string, om.OrderedMap[string, int]], om.OrderedMap[string, om.OrderedMap[string, Points]]) {
	data := om.NewOrderedMap[string, om.OrderedMap[string, int]]()
	points := om.NewOrderedMap[string, om.OrderedMap[string, Points]]()

	for id, v1 := range r1 {
		d, p := SingleFilterPass(frame.Filter(FilterEq(c1, v1)), r2, c2)

		data.Set(r1[id], d)
		points.Set(r1[id], p)
	}

	return *data, *points
}

func ProcessData(path string) om.OrderedMap[string, template.HTML] {
	frame := PrepareDataFrame(ReadDataFile(path))

	count := frame.Nrow()

	years := GetUniqueInts(frame, "year")

	span := strings.Join([]string{years[0], "-", years[len(years)-1]}, "")

	months := GetUniqueInts(frame, "month")
	months_names := *om.NewOrderedMap[string, string]()
	for _, m := range months {
		months_names.Set(m, humanmonths[m])
	}

	types := GetUniqueInts(frame, "type")
	types_names := *om.NewOrderedMap[string, string]()
	for _, t := range types {
		types_names.Set(t, frame.Filter(FilterEq("type", t)).Col("name").Records()[0])
	}

	count_years_total, points1 := SingleFilterPass(frame, years, "year")
	count_total, _ := DoubleFilterPass(frame, months, years, "month", "year")

	types_count_total, _ := SingleFilterPass(frame, types, "type")
	types_total, _ := DoubleFilterPass(frame, types, years, "type", "year")

	map_total := GeoChart("", points1)
	bar_count_total := BarChart(strings.Join([]string{"Число за", span, "(", IntToStr(count), ")"}, " "), count_years_total)
	pie_types_total := PieChart(strings.Join([]string{"Отношение за", span}, " "), SwitchKeys(types_count_total, types_names))
	bar_count_span_total := BarChartNestedValues(strings.Join([]string{"Распределение за", span}, " "), years, SwitchKeys(count_total, months_names))
	bar_types_span_total := BarChartNestedValues(strings.Join([]string{"Распределение за", span}, " "), years, SwitchKeys(types_total, types_names))

	charts_total := om.NewOrderedMap[string, template.HTML]()
	charts_total.Set("map_total", ToSnippet(map_total))
	charts_total.Set("bar_count_total", ToSnippet(bar_count_total))
	charts_total.Set("pie_types_total", ToSnippet(pie_types_total))
	charts_total.Set("bar_count_span_total", ToSnippet(bar_count_span_total))
	charts_total.Set("bar_types_span_total", ToSnippet(bar_types_span_total))

	return *charts_total
}
