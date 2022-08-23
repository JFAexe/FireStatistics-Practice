package main

import (
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

func ProcessData(path string) {
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
	count_years, _ := DoubleFilterPass(frame, years, months, "year", "month")

	types_count_total, _ := SingleFilterPass(frame, types, "type")
	types_total, points2 := DoubleFilterPass(frame, types, years, "type", "year")
	types_years, _ := DoubleFilterPass(frame, years, types, "year", "type")

	out := MakePage(
		GeoChart("all", points1),
		BarChart(strings.Join([]string{"Число за", span, "(", IntToStr(count), ")"}, " "), count_years_total),
		PieChart(strings.Join([]string{"Отношение за", span}, " "), SwitchKeys(types_count_total, types_names)),
		BarChartNestedValues(strings.Join([]string{"Распределение за", span}, " "), years, SwitchKeys(count_total, months_names)),
		BarChartNestedValues(strings.Join([]string{"Распределение за", span}, " "), years, SwitchKeys(types_total, types_names)),
	)
	for el := count_years.Front(); el != nil; el = el.Next() {
		out.AddCharts(BarChart(strings.Join([]string{"Число за", el.Key}, " "), SwitchKeys(el.Value, months_names)))
	}
	for el := types_years.Front(); el != nil; el = el.Next() {
		out.AddCharts(PieChart(strings.Join([]string{"Отношение за", el.Key}, " "), SwitchKeys(el.Value, types_names)))
	}
	points2 = SwitchKeys(points2, types_names)
	for el := points2.Front(); el != nil; el = el.Next() {
		out.AddCharts(GeoChart(el.Key, el.Value))
	}
	RenderPage(out, GetFileNameFromPath(path), "page.html")
}
