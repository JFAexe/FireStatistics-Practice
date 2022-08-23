package main

import (
	"os"
)

func init() {
	SetupLogger()
}

func main() {
	args := os.Args[1:]

	argscount := len(args)

	if argscount < 1 {
		InfoLogger.Println("No inputs provided")

		return
	}

	InfoLogger.Printf("Inputs count: %v", argscount)

	for i := 0; i < argscount; i++ {
		path := args[i]

		exists, err := IsValidFile(path)
		if err != nil {
			ErrorLogger.Printf("Something is wrong with the file \"%v\". Error: %v\n", path, err)
		}

		if !exists {
			InfoLogger.Printf("\"%v\" isn't a file or doesn't exist\n", path)

			continue
		}

		file := GetFileNameFromPath(path)

		InfoLogger.Printf("Current file: %v (%v)\n", file, path)

		ProcessData(path)

		LogMemoryUsage()
	}
}
