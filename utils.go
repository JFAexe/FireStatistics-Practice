package main

import (
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
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
	if err := os.MkdirAll(filepath.Join(temppath, path, "/"), 0700); err != nil {
		ErrorLogger.Panicf("Can't create directory. Error: %s\n", err)
	}

	file, err := os.Create(filepath.Join(temppath, path, name, "/"))
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

func RemoveDuplicateValues[T comparable](s []T) []T {
	keys := make(map[T]bool)
	list := make([]T, 0)

	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func Approx(t, a, b float64) bool {
	return math.Abs(a-b) <= t
}

func PointInSlice(r float64, p Point, s Points) bool {
	for _, c := range s {
		if Approx(r, p.x, c.x) && Approx(r, p.y, c.y) {
			return true
		}
	}

	return false
}

func FilterPoints(r float64, p Points) Points {
	sort.Slice(p, func(i, j int) bool {
		return (p[i].x > p[j].x) && (p[i].y > p[j].y)
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
	ret, err := time.Parse(timeformat, in[0])
	if err != nil {
		ErrorLogger.Panicf("Can't parse date. Error: %s\n", err)
	}

	return ret
}

func DateYear(in []string) int {
	return ParseDate(in).Year()
}

func DateMonth(in []string) int {
	return int(ParseDate(in).Month())
}

func DateDay(in []string) int {
	return ParseDate(in).Day()
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
