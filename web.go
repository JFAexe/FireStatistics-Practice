package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

const serverport string = ":1337"

func MakeHttpHandleFunc(path, text string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("got", path, "request")
		io.WriteString(w, text)
	}
}

func MakeHttpHandle(path, text string) {
	path = strings.Join([]string{"", path}, "/")

	http.HandleFunc(path, MakeHttpHandleFunc(path, text))
}

func RunHttpServer() {
	if err := http.ListenAndServe(serverport, nil); err != nil {
		panic(err)
	}
}
