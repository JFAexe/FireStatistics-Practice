package main

import (
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	om "github.com/elliotchance/orderedmap/v2"
)

const timeformat string = "2006-01-02"

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
	inf, err := os.Stat(path)
	if err == nil {
		return !inf.IsDir() && filepath.Ext(path) == ".csv", nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
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
		if _, val := keys[entry]; !val {
			keys[entry] = true

			list = append(list, entry)
		}
	}

	return list
}

func SwitchKeys[T any](m om.OrderedMap[string, T], k om.OrderedMap[string, string]) om.OrderedMap[string, T] {
	for _, oldkey := range m.Keys() {
		var newkey string

		for _, t := range k.Keys() {
			if t != oldkey {
				continue
			}

			newkey, _ = k.Get(t)

			break
		}

		val, _ := m.Get(oldkey)

		m.Delete(oldkey)
		m.Set(newkey, val)
	}

	return m
}

func InRadius(p, c Point, r float64) bool {
	x := math.Abs(c.x - p.x)
	y := math.Abs(c.y-p.y) * 1.25

	return x+y <= r
}

func SimilarInSlice(r float64, p Point, s Points) (Point, bool) {
	for _, c := range s {
		if InRadius(p, c, r) {
			return c, true
		}
	}

	return Point{}, false
}

func FilterPoints(r float64, p Points) map[Point]int {
	sort.Slice(p, func(i, j int) bool {
		return (p[i].x < p[j].x) && (p[i].y < p[j].y)
	})

	ret := make(map[Point]int, 0)
	tmp := make(Points, 0)

	for _, point := range p {
		if p, in := SimilarInSlice(r, point, tmp); in {
			ret[p]++

			continue
		}

		tmp = append(tmp, point)
		ret[point]++
	}

	return ret
}

func ParseDate(in []string) time.Time {
	ret, err := time.Parse(timeformat, in[0])
	if err != nil {
		ErrorLogger.Fatalf("Can't parse date. Error: %s\n", err)
	}

	return ret
}

func DateYear(in []string) int {
	return ParseDate(in).Year()
}

func DateMonth(in []string) int {
	return int(ParseDate(in).Month())
}

func IntToStr(in int) string {
	return strconv.Itoa(in)
}

func OpenUrlInBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)

	return exec.Command(cmd, args...).Start()
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

// Jonathan Blow is right.
