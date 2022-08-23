package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Println("No files provided")

		return
	}

	for i := 0; i < len(args); i++ {
		path := args[i]

		exists, err := IsValidFile(path)
		if err != nil {
			log.Println(err)
		}

		if !exists {
			log.Println(path, "isn't a file or doesn't exist")

			continue
		}

		file := GetFileNameFromPath(path)

		log.Println("FILE:", file)

		ProcessData(path)
	}
}
