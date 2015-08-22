package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
)

func handleNow(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, time.Now())
}

func main() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" {
		port = "8000"
	}

	n := negroni.Classic()
	n.UseHandler(http.HandlerFunc(handleNow))
	n.Run(host + ":" + port)
}
