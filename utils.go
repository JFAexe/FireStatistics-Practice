package main

import (
	"log"
	"math"
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
		return !info.IsDir() && filepath.Ext(path) == ".csv", nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CreateFile(path, name string) *os.File {
	if err := os.MkdirAll(strings.Join([]string{temppath, path}, "/"), 0700); err != nil {
		ErrorLogger.Panicf("Can't create directory. Error: %s\n", err)
	}

	file, err := os.Create(strings.Join([]string{temppath, path, name}, "/"))
	if err != nil {
		ErrorLogger.Panicf("Can't create file. Error: %s\n", err)
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

func RemoveDuplicateStrings(slice []string) []string {
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

func Approx(t, a, b float64) bool {
	return math.Abs(a-b) <= t
}

func PointInSlice(r float64, p []float64, s Points) bool {
	for _, c := range s {
		if Approx(r, p[0], c[0]) && Approx(r, p[1], c[1]) {
			return true
		}
	}

	return false
}

func FilterPoints(r float64, p Points) Points {
	sort.Slice(p, func(i, j int) bool {
		return (p[i][0] > p[j][0]) && (p[i][1] > p[j][1])
	})

	temp := make(Points, 0)

	for _, point := range p {
		if PointInSlice(r, point, temp) {
			continue
		}

		temp = append(temp, point)
	}

	return temp
}

func ParseDate(in []string) time.Time {
	ret, err := time.Parse(timeformat, strings.Join(in, ""))
	if err != nil {
		ErrorLogger.Panicf("Can't parse date. Error: %s\n", err)
	}

	return ret
}

func ParseNumber(in string) int {
	ret, err := strconv.Atoi(in)
	if err != nil {
		ErrorLogger.Panicf("Can't parse number. Error: %s\n", err)
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
