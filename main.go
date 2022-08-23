package main

import "os"

func init() {
	SetupLogger()
}

func main() {
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
			ErrorLogger.Printf("Something is wrong with the file \"%v\". Error: %v\n", arg, err)
		}

		if !exists {
			InfoLogger.Printf("\"%v\" isn't a file or doesn't exist\n", arg)

			continue
		}

		InfoLogger.Printf("Current file: %v (%v)\n", GetFileNameFromPath(arg), arg)

		ProcessData(arg)

		LogMemoryUsage()
	}
}
