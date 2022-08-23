package main

import (
	"fmt"
	"os"

	dataframe "github.com/go-gota/gota/dataframe"
	series "github.com/go-gota/gota/series"
)

func ReadFile(path string) dataframe.DataFrame {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return dataframe.ReadCSV(file, dataframe.WithDelimiter(';'))
}

func GetSortedDF(df dataframe.DataFrame, s string) dataframe.DataFrame {
	return df.Copy().Arrange(dataframe.Sort(s))
}

func ProcessData(path string) {
	df := ReadFile(path)

	dates := df.Select("dt").Records()[1:]

	df = df.Mutate(series.New(Map(dates, func(i []string) int { return ParseDate(i).Year() }), series.Int, "year"))
	df = df.Mutate(series.New(Map(dates, func(i []string) int { return int(ParseDate(i).Month()) }), series.Int, "month"))
	df = df.Mutate(series.New(Map(dates, func(i []string) int { return ParseDate(i).Day() }), series.Int, "day"))

	fmt.Println(df)
}
