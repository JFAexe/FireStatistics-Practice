package main

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

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
	val, _ := omap.Get(key)
	return val
}

func RunHTTPServer(server *http.Server) {
	InfoLogger.Printf("Starting server on %s\n", Port)

	assets, err := fs.Sub(embedded, "assets")
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets))))

	tplfuncs := template.FuncMap{
		"inc":           func(n int) int { return n + 1 },
		"safe":          func(s string) template.HTML { return template.HTML(s) },
		"getkeysHTML":   TplGetKeys[template.HTML],
		"getbykeyHTML":  TplGetByKey[template.HTML],
		"getkeysBlock":  TplGetKeys[Block],
		"getbykeyBlock": TplGetByKey[Block],
	}

	HTMLTemplates = template.Must(template.New("templates.html").Funcs(tplfuncs).ParseFS(assets, "templates.html"))

	if err := server.ListenAndServe(); err != nil {
		switch err {
		case http.ErrServerClosed:
			InfoLogger.Println("Server shut down")
		default:
			ErrorLogger.Fatal(err)
		}
	}
}

func ShutdownHTTPServer(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit

	InfoLogger.Println("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		ErrorLogger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	close(done)
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

// git -gud
