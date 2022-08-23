package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		path := args[i]

		exists, err := IsValidFile(path)
		if err != nil {
			fmt.Println(err)
		}

		if !exists {
			fmt.Println(path, "isn't a file or doesn't exist")

			continue
		}

		fmt.Println("FILE:", path)

		ProcessData(path)
	}
}
