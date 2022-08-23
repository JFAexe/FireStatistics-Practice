package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	temppath   string = "./fa-temp"
	timeformat string = "2006-01-02"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func SetupLogger() {
	var logout *os.File

	logout = os.Stdout

	InfoLogger = log.New(logout, "[INFO] ", log.LstdFlags|log.Lmsgprefix)
	ErrorLogger = log.New(logout, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
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

func WriteFileFromBytes(path, name string, buf []byte) error {
	dir := strings.Join([]string{temppath, GetFileNameFromPath(path)}, "/")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(dir, name), buf, 0600); err != nil {
		return err
	}

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

	prev := 1

	for cur := 1; cur < len(s); cur++ {
		if s[cur-1] != s[cur] {
			s[prev] = s[cur]
			prev++
		}
	}

	return s[:prev]
}

func GetFileNameFromPath(path string) string {
	_, file := filepath.Split(path)

	return strings.TrimSuffix(file, filepath.Ext(file))
}

func ParseDate(i []string) time.Time {
	date, err := time.Parse(timeformat, strings.Join(i, ""))
	if err != nil {
		ErrorLogger.Panicf("Can't parse date. Error: %v\n", err)
	}

	return date
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
