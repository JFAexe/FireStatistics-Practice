package main

import (
	"io"
	"net/http"
	"strings"
)

const serverport string = ":1337"

func MakeHttpHandle(path, text string) {
	path = strings.Join([]string{"", path}, "/")

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		InfoLogger.Println("got", path, "request")
		io.WriteString(w, text)
	})
}

func RunHttpServer() {
	if err := http.ListenAndServe(serverport, nil); err != nil {
		ErrorLogger.Panic(err)
	}
}
