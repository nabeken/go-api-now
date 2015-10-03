package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
)

func handleEnvNow(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, time.Now())

	for _, env := range os.Environ() {
		fmt.Fprintln(w, env)
	}
}

func main() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" {
		port = "8000"
	}

	n := negroni.Classic()
	n.UseHandler(http.HandlerFunc(handleEnvNow))
	n.Run(host + ":" + port)
}
