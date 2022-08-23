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

const (
	DotRad      float64 = 4
	DotDia      float32 = float32(DotRad * 18)
	HTMLSnippet string  = `
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
)

var SnippetTemplate *template.Template

func NewSnippetRenderer(c interface{}, before ...func()) Renderer {
	return &SnippetRenderer{c: c, before: before}
}

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

func WithRenderer(opt Renderer) ch.GlobalOpts {
	return func(bc *ch.BaseConfiguration) { bc.Renderer = opt }
}

func WithScatterSize(opt float32) ch.SeriesOpts {
	return func(s *ch.SingleSeries) { s.SymbolSize = opt }
}

func DefaultOptions(title string, chart interface{ Validate() }) []ch.GlobalOpts {
	return []ch.GlobalOpts{
		ch.WithTitleOpts(op.Title{Title: title, Left: "center", TitleStyle: &op.TextStyle{FontFamily: "'Exo 2', sans-serif"}}),
		ch.WithInitializationOpts(op.Initialization{Width: "1080px", Height: "400px"}),
		ch.WithTooltipOpts(op.Tooltip{Show: true}),
		ch.WithLegendOpts(op.Legend{Show: true, Bottom: "bottom", Left: "center"}),
		WithRenderer(NewSnippetRenderer(chart, chart.Validate)),
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

func BarChart(title string, values om.OrderedMap[string, int]) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	series, axis := ConverDataBar(values)

	chart.AddSeries("", series)

	chart.SetXAxis(axis).SetSeriesOptions(ch.WithLabelOpts(op.Label{Show: true, Position: "top"}))

	return chart
}

func BarChartNestedValues(title string, axis []string, values om.OrderedMap[string, om.OrderedMap[string, int]]) *ch.Bar {
	chart := ch.NewBar()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	for _, key := range values.Keys() {
		data, _ := values.Get(key)
		series, _ := ConverDataBar(data)
		chart.AddSeries(key, series)
	}

	chart.SetXAxis(axis)

	return chart
}

func ConverDataPie(data om.OrderedMap[string, int]) []op.PieData {
	ret := make([]op.PieData, data.Len())

	for id, key := range data.Keys() {
		value, _ := data.Get(key)
		ret[id] = op.PieData{Name: key, Value: value}
	}

	return ret
}

func PieChart(title string, values om.OrderedMap[string, int]) *ch.Pie {
	chart := ch.NewPie()

	chart.SetGlobalOptions(DefaultOptions(title, chart)...)

	chart.AddSeries("", ConverDataPie(values))

	chart.SetSeriesOptions(ch.WithPieChartOpts(op.PieChart{Radius: []string{"25%", "55%"}}))

	return chart
}

func ConverDataGeo(data Points, tip string) ([]op.GeoData, float32, float32) {
	ret := make([]op.GeoData, 0)

	min, max := 1, 1

	for point, count := range FilterPoints(DotRad, data) {
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

func GeoChart(title string, data om.OrderedMap[string, Points]) *ch.Geo {
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

		chart.AddSeries(key, tp.ChartScatter, series, WithScatterSize(DotDia))
	}

	chart.SetGlobalOptions(append(
		DefaultOptions(title, chart),
		ch.WithGeoComponentOpts(op.GeoComponent{Map: "Russia"}),
		ch.WithLegendOpts(op.Legend{Show: true, Bottom: "bottom", Left: "center"}),
		ch.WithVisualMapOpts(op.VisualMap{Calculable: true, Min: cmin, Max: cmax}),
	)...)

	return chart
}
