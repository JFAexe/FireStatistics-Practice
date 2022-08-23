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
	Filter func(string, string) df.F
)

var humanmonths = map[string]string{
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

func ReadDataFile(path string) df.DataFrame {
	file, err := os.Open(path)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
	defer file.Close()

	return df.ReadCSV(file, df.WithDelimiter(';'))
}

func PrepareDataFrame(fr df.DataFrame) df.DataFrame {
	dt := fr.Select("dt").Records()[1:]

	return fr.
		Drop("dt").
		Rename("type", "type_id").
		Rename("name", "type_name").
		Mutate(sr.New(Map(dt, DateYear), sr.Int, "year")).
		Mutate(sr.New(Map(dt, DateMonth), sr.Int, "month")).
		Arrange(df.Sort("year"))
}

func FilterEq(col, val string) df.F {
	return df.F{Colname: col, Comparator: sr.Eq, Comparando: val}
}

func GetUniqueInts(fr df.DataFrame, name string) []string {
	ret, err := fr.Col(name).Int()
	if err != nil {
		ErrorLogger.Fatalf("Can't parse int. Error: %s\n", err)
	}

	ret = RemoveDuplicateValues(ret)

	sort.Ints(ret)

	return Map(ret, IntToStr)
}

func GetPoints(fr df.DataFrame) Points {
	ret := make(Points, 0)

	lon := fr.Col("lon").Float()
	lat := fr.Col("lat").Float()

	for id := 0; id < len(lon); id++ {
		ret = append(ret, Point{lon[id], lat[id]})
	}

	return ret
}

func SingleFilterPass(fr df.DataFrame, r []string, n string, f Filter) (om.OrderedMap[string, int], om.OrderedMap[string, Points]) {
	data := om.NewOrderedMap[string, int]()
	points := om.NewOrderedMap[string, Points]()

	for id, v := range r {
		f := fr.Filter(f(n, v))

		data.Set(r[id], f.Nrow())
		points.Set(r[id], GetPoints(f))
	}

	return *data, *points
}

func DoubleFilterPass(fr df.DataFrame, r1, r2 []string, n1, n2 string, f Filter) (om.OrderedMap[string, om.OrderedMap[string, int]], om.OrderedMap[string, om.OrderedMap[string, Points]]) {
	data := om.NewOrderedMap[string, om.OrderedMap[string, int]]()
	points := om.NewOrderedMap[string, om.OrderedMap[string, Points]]()

	for id, v1 := range r1 {
		d, p := SingleFilterPass(fr.Filter(f(n1, v1)), r2, n2, f)

		data.Set(r1[id], d)
		points.Set(r1[id], p)
	}

	return *data, *points
}

func ProcessData(path string) Page {
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

	names := map[string]string{
		"bar_count":           strings.Join([]string{"Подсчёт случаев за ", span}, ""),
		"bar_count_title":     strings.Join([]string{"Суммарно зарегистрировано ", IntToStr(count)}, ""),
		"bar_span":            strings.Join([]string{"Распределение за", span}, " "),
		"bar_year_count":      "Подсчёт случаев",
		"pie_percentage":      strings.Join([]string{"Доли типов за", span}, " "),
		"pie_year_percentage": "Доли типов",
		"map":                 "Карта распределения",
		"map_year_count":      "Распределение",
		"map_year_types":      "Карта типов",
	}

	counts_years_total, points_counts_years_total := SingleFilterPass(frame, years, "year", FilterEq)
	counts_total, points_counts_total := DoubleFilterPass(frame, months, years, "month", "year", FilterEq)
	counts_years, points_counts_years := DoubleFilterPass(frame, years, months, "year", "month", FilterEq)

	types_count_total, points_types_count_total := SingleFilterPass(frame, types, "type", FilterEq)
	types_total, points_types_total := DoubleFilterPass(frame, types, years, "type", "year", FilterEq)
	types_years, points_types_years := DoubleFilterPass(frame, years, types, "year", "type", FilterEq)

	charts_counts_total := *om.NewOrderedMap[string, template.HTML]()
	charts_counts_total.Set(names["map"], GeoChart("", points_counts_years_total))
	charts_counts_total.Set(names["bar_count"], BarChart(names["bar_count_title"], counts_years_total))
	charts_counts_total.Set(names["bar_span"], BarChartSeveral("", years, SwitchKeys(counts_total, months_names)))

	charts_types_total := *om.NewOrderedMap[string, template.HTML]()
	charts_types_total.Set(names["map"], GeoChart("", SwitchKeys(points_types_count_total, types_names)))
	charts_types_total.Set(names["pie_percentage"], PieChart("", SwitchKeys(types_count_total, types_names)))
	charts_types_total.Set(names["bar_span"], BarChartSeveral("", years, SwitchKeys(types_total, types_names)))

	maps_count := GeoChartNested(points_counts_total, months_names)
	maps_types := GeoChartNested(points_types_total, types_names)

	primary := []Block{
		{Id: "charts_count_total", Header: "Количество", Snippets: charts_counts_total},
		{Id: "maps_count", Snippets: maps_count},
		{Id: "charts_types_total", Header: "Типы", Snippets: charts_types_total},
		{Id: "maps_types", Snippets: maps_types},
	}

	secondary := make([]Block, len(years))
	for id, key := range counts_years.Keys() {
		charts_year := *om.NewOrderedMap[string, template.HTML]()

		mapcount, _ := points_counts_years.Get(key)
		charts_year.Set(names["map_year_count"], GeoChart("", SwitchKeys(mapcount, months_names)))

		barvalues, _ := counts_years.Get(key)
		charts_year.Set(names["bar_year_count"], BarChart("", SwitchKeys(barvalues, months_names)))

		maptypes, _ := points_types_years.Get(key)
		charts_year.Set(names["map_year_types"], GeoChart("", SwitchKeys(maptypes, types_names)))

		pievalues, _ := types_years.Get(key)
		charts_year.Set(names["pie_year_percentage"], PieChart("", SwitchKeys(pievalues, types_names)))

		secondary[id] = Block{
			Id:       key,
			Header:   key,
			Snippets: charts_year,
		}
	}

	page := Page{
		OldMap:          UseOldMap,
		Header:          strings.Join([]string{"FSP | ", GetFileNameFromPath(path)}, ""),
		ChartsPrimary:   primary,
		ChartsSecondary: secondary,
	}

	return page
}

// JFAexe was here
