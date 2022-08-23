package main

import (
	"os"
	"strings"
	"time"
)

const format string = "2006-01-02"

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

func ParseDate(i []string) time.Time {
	date, err := time.Parse(format, strings.Join(i, ""))
	if err != nil {
		panic(err)
	}

	return date
}

func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func Remove[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}
