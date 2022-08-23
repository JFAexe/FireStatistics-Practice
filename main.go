package main

import (
	"os"
	"time"
)

const (
	rad float64 = 4
	dia float32 = float32(rad * 18)
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

		InfoLogger.Printf("Current file: %s (%s)\n", GetFileNameFromPath(arg), arg)

		ProcessData(arg)

		LogMemoryUsage()
	}

	InfoLogger.Println("Main", time.Since(startup))
}
