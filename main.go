package main

import (
	"flag"
	"strings"
	"time"
)

func init() {
	SetupLogger()

	flag.BoolVar(&UseOldMap, "oldmap", false, "Use old map type for charts (No Crimea)")
	flag.StringVar(&Port, "port", ":1337", "Localhost custom port")
	flag.Float64Var(&PointRad, "pointradius", 2, "Radius of area around point on map to optimize data")

	flag.Parse()

	PointDia = float32(PointRad * 18)
}

func main() {
	startup := time.Now()

	args := flag.Args()

	if count := len(args); count < 1 {
		InfoLogger.Println("No inputs provided")

		return
	} else {
		InfoLogger.Printf("Inputs count: %v", count)
	}

	links := make(map[string]string, 0)
	link := strings.Join([]string{"http://localhost", Port}, "")

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

		page := ProcessData(arg)

		AddPageHandle(name, "pagereport", page)

		links[arg] = strings.Join([]string{link, name}, "/")
	}

	if len(links) < 1 {
		return
	}

	AddPageHandle("", "pagemain", links)

	OpenUrlInBrowser(link)

	InfoLogger.Println("Main", time.Since(startup))

	LogMemoryUsage()

	RunHTTPServer()
}
