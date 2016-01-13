package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.GET("/", gin.WrapF(handleEnvNow))
	r.Run(host + ":" + port)
}
