package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	i, err := os.Stat(path)
	if err == nil {
		return !i.IsDir(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func WriteFileFromBytes(path, name string, buf []byte) error {
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(path, name), buf, 0600); err != nil {
		return err
	}

	InfoLogger.Printf("Saved file \"%v\" at %v", name, path)

	return nil
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

func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))

	for i, v := range s {
		r[i] = f(v)
	}

	return r
}

func RemoveDuplicates(s []string) []string {
	if len(s) < 1 {
		return s
	}

	sort.Strings(s)

	prv := 1
	for cur := 1; cur < len(s); cur++ {
		if s[cur-1] != s[cur] {
			s[prv] = s[cur]
			prv++
		}
	}

	return s[:prv]
}

func GetFileNameFromPath(path string) string {
	_, f := filepath.Split(path)

	return strings.TrimSuffix(f, filepath.Ext(f))
}

func ParseDate(i []string) time.Time {
	d, err := time.Parse(timeformat, strings.Join(i, ""))
	if err != nil {
		ErrorLogger.Panicf("Can't parse date. Error: %v\n", err)
	}

	return d
}

func ParseNumber(i string) int {
	n, err := strconv.Atoi(i)
	if err != nil {
		ErrorLogger.Panicf("Can't parse number. Error: %v\n", err)
	}

	return n
}

func LogMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	InfoLogger.Printf("Alloc: %v MiB | TotalAlloc: %v MiB | Sys: %v MiB | NumGC: %v\n",
		bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
