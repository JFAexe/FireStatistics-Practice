package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	om "github.com/elliotchance/orderedmap/v2"
)

type (
	Block struct {
		Id       string
		Header   string
		Snippets om.OrderedMap[string, template.HTML]
	}

	Page struct {
		OldMap          bool
		Header          string
		ChartsPrimary   []Block
		ChartsSecondary []Block
	}
)

var (
	//go:embed assets
	embedded embed.FS

	HTMLTemplates *template.Template

	UseOldMap bool
	Port      string
)

func TplGetKeys[T any](omap om.OrderedMap[string, T]) []string {
	return omap.Keys()
}

func TplGetByKey[T any](omap om.OrderedMap[string, T], key string) T {
	value, _ := omap.Get(key)

	return value
}

func RunHTTPServer() {
	assets, err := fs.Sub(embedded, "assets")
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	fileserver := http.FileServer(http.FS(assets))
	fileserver = http.StripPrefix("/assets/", fileserver)

	http.Handle("/assets/", fileserver)

	tplfuncs := template.FuncMap{
		"inc":           func(num int) int { return num + 1 },
		"safe":          func(s string) template.HTML { return template.HTML(s) },
		"getkeysHTML":   TplGetKeys[template.HTML],
		"getbykeyHTML":  TplGetByKey[template.HTML],
		"getkeysBlock":  TplGetKeys[Block],
		"getbykeyBlock": TplGetByKey[Block],
	}

	HTMLTemplates = template.Must(template.New("templates.html").Funcs(tplfuncs).ParseFS(assets, "templates.html"))

	if err := http.ListenAndServe(Port, nil); err != nil {
		ErrorLogger.Fatal(err)
	}
}

func AddPageHandle(path, tmpl string, data any) {
	path = strings.Join([]string{"", path}, "/")

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		InfoLogger.Println("Got", path, "request")

		if err := HTMLTemplates.ExecuteTemplate(w, tmpl, data); err != nil {
			ErrorLogger.Fatal(err)
		}
	})
}
