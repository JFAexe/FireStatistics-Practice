package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	temppath   string = "./temp"
	timeformat string = "2006-01-02"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func SetupLogger() {
	out := os.Stdout
	flags := log.LstdFlags | log.Lmsgprefix

	InfoLogger = log.New(out, "[INFO] ", flags)
	ErrorLogger = log.New(out, "[ERROR] ", flags)
}

func IsValidFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CreateFile(path, name string) *os.File {
	if err := os.MkdirAll(strings.Join([]string{temppath, path}, "/"), 0700); err != nil {
		ErrorLogger.Panicf("Can't create directory. Error: %v\n", err)
	}

	file, err := os.Create(strings.Join([]string{temppath, path, name}, "/"))
	if err != nil {
		ErrorLogger.Panicf("Can't create file. Error: %v\n", err)
	}

	return file
}

func GetFileNameFromPath(path string) string {
	_, file := filepath.Split(path)

	return strings.TrimSuffix(file, filepath.Ext(file))
}

func Map[T, U any](slice []T, fn func(T) U) []U {
	ret := make([]U, len(slice))

	for id, val := range slice {
		ret[id] = fn(val)
	}

	return ret
}

func RemoveDuplicates(slice []string) []string {
	if len(slice) < 1 {
		return slice
	}

	sort.Strings(slice)

	prev := 1
	for curr := 1; curr < len(slice); curr++ {
		if slice[curr-1] != slice[curr] {
			slice[prev] = slice[curr]
			prev++
		}
	}

	return slice[:prev]
}

func ParseDate(in []string) time.Time {
	ret, err := time.Parse(timeformat, strings.Join(in, ""))
	if err != nil {
		ErrorLogger.Panicf("Can't parse date. Error: %v\n", err)
	}

	return ret
}

func ParseNumber(in string) int {
	ret, err := strconv.Atoi(in)
	if err != nil {
		ErrorLogger.Panicf("Can't parse number. Error: %v\n", err)
	}

	return ret
}

func ParseFloatArray(in []string) []float64 {
	ret := make([]float64, 0)

	for _, val := range in {
		float, err := strconv.ParseFloat(val, 64)
		if err != nil {
			ErrorLogger.Panicf("Can't parse float. Error: %v\n", err)
		}

		ret = append(ret, float)
	}

	return ret
}

func LogMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	InfoLogger.Printf("Alloc: %v MiB | TotalAlloc: %v MiB | Sys: %v MiB | NumGC: %v\n",
		bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1048576
}
