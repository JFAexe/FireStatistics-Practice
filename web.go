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
		Header  string
		Blocks  om.OrderedMap[string, Block]
		Tabs    Block
		Comment template.HTML
	}
)

const (
	serverport string        = ":1337"
	htmlmsg    template.HTML = template.HTML(`<!-- Welcome to the cursed world we created by ourselves -->`)
)

var (
	//go:embed assets
	embedded      embed.FS
	HTMLTemplates *template.Template
)

func RunHttpServer() {
	assets, err := fs.Sub(embedded, "assets")
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	fileserver := http.FileServer(http.FS(assets))
	fileserver = http.StripPrefix("/assets/", fileserver)

	http.Handle("/assets/", fileserver)

	tplfuncs := template.FuncMap{
		"inc": func(num int) int {
			return num + 1
		},
		"getkeysHTML": func(omap om.OrderedMap[string, template.HTML]) []string {
			return omap.Keys()
		},
		"getbykeyHTML": func(omap om.OrderedMap[string, template.HTML], key string) template.HTML {
			value, _ := omap.Get(key)
			return value
		},
		"getkeysBlock": func(omap om.OrderedMap[string, Block]) []string {
			return omap.Keys()
		},
		"getbykeyBlock": func(omap om.OrderedMap[string, Block], key string) Block {
			value, _ := omap.Get(key)
			return value
		},
	}

	HTMLTemplates = template.Must(template.New("templates.html").Funcs(tplfuncs).ParseFiles("assets/templates.html"))

	if err := http.ListenAndServe(serverport, nil); err != nil {
		ErrorLogger.Fatal(err)
	}
}

func MakeHttpHandle(path, tmpl string, data Page) {
	path = strings.Join([]string{"", path}, "/")

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		InfoLogger.Println("got", path, "request")

		if err := HTMLTemplates.ExecuteTemplate(w, tmpl, data); err != nil {
			ErrorLogger.Fatal(err)
		}
	})
}
