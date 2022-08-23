package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	ChanDone = make(chan bool, 1)
	ChanQuit = make(chan os.Signal, 1)
)

func init() {
	SetupLogger()

	signal.Notify(ChanQuit, os.Interrupt)

	flag.BoolVar(&UseOldMap, "oldmap", false, "Use old map type for charts (No Crimea)")
	flag.StringVar(&Port, "port", ":1337", "Localhost custom port")
	flag.Float64Var(&PointRad, "radius", 2, "Radius of area around point on map to optimize data")

	flag.Parse()

	PointDia = float32(PointRad * 18)
}

func main() {
	args := flag.Args()

	if count := len(args); count < 1 {
		InfoLogger.Println("No inputs provided")

		return
	} else {
		InfoLogger.Printf("Inputs count: %v", count)
	}

	server := &http.Server{Addr: Port}

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

		begin := time.Now()

		AddPageHandle(name, "pagereport", ProcessData(arg))

		links[arg] = strings.Join([]string{link, name}, "/")

		InfoLogger.Printf("File %s took %v", name, time.Since(begin))
	}

	if len(links) < 1 {
		return
	}

	AddPageHandle("", "pagemain", links)

	go RunHTTPServer(server)

	go ShutdownHTTPServer(server, ChanQuit, ChanDone)

	OpenUrlInBrowser(link)

	<-ChanDone

	LogMemoryUsage()
}

// Daft is dead, only punk remains
