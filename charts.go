package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"regexp"

	om "github.com/elliotchance/orderedmap/v2"
	ch "github.com/go-echarts/go-echarts/v2/charts"
	op "github.com/go-echarts/go-echarts/v2/opts"
	tp "github.com/go-echarts/go-echarts/v2/types"
)

type (
	Renderer interface {
		Render(w io.Writer) error
	}

	SnippetRenderer struct {
		c      interface{}
		before []func()
	}
)

const HTMLSnippet string = `
<div class = 'charts-container'>
    <div class = 'item' id = '{{ .ChartID }}' style = 'width:{{ .Initialization.Width }}; height:{{ .Initialization.Height }};'></div>
</div>
<script type = 'text/javascript'> 'use strict';
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), '{{ .Theme }}');
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    	{{ . | safeJS }}
    {{- end }}
</script>
`

var (
	PointRad float64
	PointDia float32
)

func (r *SnippetRenderer) Render(w io.Writer) error {
	for _, fn := range r.before {
		fn()
	}

	tplfuncs := template.FuncMap{
		"safeJS": func(s interface{}) template.JS {
			return template.JS(fmt.Sprint(s))
		},
	}

	tpl := template.Must(template.New("chartsnippet").Funcs(tplfuncs).Parse(HTMLSnippet))

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "chartsnippet", r.c); err != nil {
		return err
	}

	pat := regexp.MustCompile(`(__f__")|("__f__)|(__f__)`)
	content := pat.ReplaceAll(buf.Bytes(), []byte(""))

	_, err := w.Write(content)

	return err
}

func ToSnippet(r Renderer) template.HTML {
	var buf bytes.Buffer

	r.Render(&buf)

	return template.HTML(buf.String())
}

func WithSnippetRenderer(chart interface{ Validate() }) ch.GlobalOpts {
	return func(bc *ch.BaseConfiguration) {
		bc.Renderer = &SnippetRenderer{c: chart, before: []func(){chart.Validate}}
	}
}

func WithScatterSize(opt float32) ch.SeriesOpts {
	return func(s *ch.SingleSeries) { s.SymbolSize = opt }
}

func GetChartOptions(title string, chart interface{ Validate() }) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title, Left: "center", TitleStyle: &op.TextStyle{FontFamily: "'Exo 2', sans-serif"}}),
		ch.WithInitializationOpts(op.Initialization{Width: "1080px", Height: "400px"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
		ch.WithLegendOpts(op.Legend{Show: true, Bottom: "bottom", Left: "center"}),
		ch.WithToolboxOpts(op.Toolbox{
			Show:   true,
			Right:  "5%",
			Top:    "center",
			Orient: "vertical",
			Feature: &op.ToolBoxFeature{
				SaveAsImage: &op.ToolBoxFeatureSaveAsImage{Show: true, Type: "png", Title: "График"},
				DataView:    &op.ToolBoxFeatureDataView{Show: true, Title: "Данные", Lang: []string{"Исходные данные", "Закрыть", "Обновить"}},
			}},
		),
		WithSnippetRenderer(chart),
	}
}

func ConverDataBar(data om.OrderedMap[string, int]) ([]op.BarData, []string) {
	ret := make([]op.BarData, data.Len())
	axs := make([]string, data.Len())

	for id, key := range data.Keys() {
		value, _ := data.Get(key)
		ret[id] = op.BarData{Name: key, Value: value}
		axs[id] = key
	}

	return ret, axs
}

func BarChart(title string, values om.OrderedMap[string, int]) template.HTML {
	chart := ch.NewBar()

	chart.SetGlobalOptions(GetChartOptions(title, chart)...)

	series, axis := ConverDataBar(values)

	chart.AddSeries("", series)

	chart.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return ToSnippet(chart)
}

func BarChartSeveral(title string, axis []string, values om.OrderedMap[string, om.OrderedMap[string, int]]) template.HTML {
	chart := ch.NewBar()

	chart.SetGlobalOptions(GetChartOptions(title, chart)...)

	for _, key := range values.Keys() {
		data, _ := values.Get(key)
		series, _ := ConverDataBar(data)
		chart.AddSeries(key, series)
	}

	chart.SetXAxis(axis)

	return ToSnippet(chart)
}

func ConverDataPie(data om.OrderedMap[string, int]) []op.PieData {
	ret := make([]op.PieData, data.Len())

	for id, key := range data.Keys() {
		value, _ := data.Get(key)
		ret[id] = op.PieData{Name: key, Value: value}
	}

	return ret
}

func PieChart(title string, values om.OrderedMap[string, int]) template.HTML {
	chart := ch.NewPie()

	chart.SetGlobalOptions(GetChartOptions(title, chart)...)

	chart.AddSeries("", ConverDataPie(values))

	chart.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return ToSnippet(chart)
}

func ConverDataGeo(data Points, tip string) ([]op.GeoData, float32, float32) {
	ret := make([]op.GeoData, 0)

	min, max := 1, 1

	for point, count := range FilterPoints(PointRad, data) {
		if count < min {
			min = count
		}

		if count > max {
			max = count
		}

		ret = append(ret, op.GeoData{Name: tip, Value: []any{point.x, point.y, count}})
	}

	return ret, float32(min), float32(max)
}

func GeoChart(title string, data om.OrderedMap[string, Points]) template.HTML {
	chart := ch.NewGeo()

	var cmin, cmax float32 = 1, 1

	for _, key := range data.Keys() {
		points, _ := data.Get(key)
		series, min, max := ConverDataGeo(points, key)

		if cmin < cmin {
			cmin = min
		}

		if max > cmax {
			cmax = max
		}

		chart.AddSeries(key, tp.ChartScatter, series, WithScatterSize(PointDia))
	}

	chart.SetGlobalOptions(append(
		GetChartOptions(title, chart),
		ch.WithGeoComponentOpts(op.GeoComponent{Map: "Russia"}),
		ch.WithVisualMapOpts(op.VisualMap{Calculable: true, Min: cmin, Max: cmax}),
	)...)

	return ToSnippet(chart)
}

func GeoChartNested(m om.OrderedMap[string, om.OrderedMap[string, Points]], k om.OrderedMap[string, string]) om.OrderedMap[string, template.HTML] {
	ret := om.NewOrderedMap[string, template.HTML]()

	m = SwitchKeys(m, k)
	for _, key := range m.Keys() {
		value, _ := m.Get(key)
		ret.Set(key, GeoChart("", value))
	}

	return *ret
}
