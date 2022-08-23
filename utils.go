package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const format string = "2006-01-02"
const temppath string = "./temp"

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

func WriteFileFromBytes(path, suffix string, buf []byte) error {
	file := GetFileNameFromPath(path)

	dir := strings.Join([]string{temppath, file}, "/")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	out := filepath.Join(dir, strings.Join([]string{file, suffix}, ""))

	if err := ioutil.WriteFile(out, buf, 0600); err != nil {
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

	for curr := 1; curr < len(s); curr++ {
		if s[curr-1] != s[curr] {
			s[prev] = s[curr]
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
	date, err := time.Parse(format, strings.Join(i, ""))
	if err != nil {
		panic(err)
	}

	return date
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
