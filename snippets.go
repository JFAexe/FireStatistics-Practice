package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	rn "github.com/go-echarts/go-echarts/v2/render"
)

const (
	tplname    string = "chartsnippet"
	tplsnippet string = `
<script src="../../../web/echarts.min.js"></script>
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
</div>
<script type="text/javascript"> "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    	{{ . | safeJS }}
    {{- end }}
</script>
`
)

type SnippetRenderer struct {
	c      interface{}
	before []func()
}

func NewSnippetRenderer(c interface{}, before ...func()) rn.Renderer {
	return &SnippetRenderer{c: c, before: before}
}

func (r *SnippetRenderer) Render(w io.Writer) error {
	for _, fn := range r.before {
		fn()
	}

	tpl := template.Must(
		template.New(tplname).
			Funcs(template.FuncMap{"safeJS": func(s interface{}) template.JS { return template.JS(fmt.Sprint(s)) }}).
			Parse(tplsnippet),
	)

	err := tpl.ExecuteTemplate(w, tplname, r.c)

	return err
}

func ChartToSnippet(r rn.Renderer) string {
	var buf bytes.Buffer

	r.Render(&buf)

	return buf.String()
}
