package main

import (
	"html/template"
	"os"
	"strings"
	"time"

	om "github.com/elliotchance/orderedmap/v2"
)

func init() {
	SetupLogger()
}

func main() {
	startup := time.Now()

	args := os.Args[1:]

	if count := len(args); count < 1 {
		InfoLogger.Println("No inputs provided")

		return
	} else {
		InfoLogger.Printf("Inputs count: %v", count)
	}

	for _, arg := range args {
		exists, err := IsValidFile(arg)
		if err != nil {
			ErrorLogger.Printf("Something is wrong with the file \"%s\". Error: %s\n", arg, err)
		}

		if !exists {
			InfoLogger.Printf("\"%s\" isn't a csv file or doesn't exist\n", arg)

			continue
		}

		name := GetFileNameFromPath(arg)

		InfoLogger.Printf("Current file: %s (%s)\n", name, arg)

		data := ProcessData(arg)

		pageblock := *om.NewOrderedMap[string, Block]()
		pageblock.Set("total", Block{
			Id:       "total",
			Header:   "Заголовок",
			Snippets: data,
		})

		tabblock := *om.NewOrderedMap[string, template.HTML]()
		tabblock.Set("1", template.HTML("<p>Test 1</p>"))
		tabblock.Set("2", template.HTML("<p>Test 2</p>"))

		page := Page{
			Header: strings.Join([]string{"FSP (", name, ")"}, ""),
			Blocks: pageblock,
			Tabs: Block{
				Id:       "test",
				Header:   "Test",
				Snippets: tabblock,
			},
			Comment: htmlmsg,
		}

		MakeHttpHandle(name, "document", page)

		OpenUrlInBrowser(strings.Join([]string{"http://localhost:1337/", name}, ""))
	}

	InfoLogger.Println("Main", time.Since(startup))

	LogMemoryUsage()

	RunHttpServer()
}
