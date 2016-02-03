package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/gin-gonic/gin"
)

var Version = "201601170013"

func printEnvNow(w io.Writer) {
	fmt.Fprintln(w, "version:", Version)
	fmt.Fprintln(w, time.Now())

	for _, env := range os.Environ() {
		fmt.Fprintln(w, env)
	}
}

func main() {
	var httpMode = flag.Bool("http", true, "enable HTTP server")
	flag.Parse()

	if *httpMode {
		HTTP()
	} else {
		for range time.Tick(10 * time.Second) {
			printEnvNow(os.Stdout)
		}
	}
}

func HTTP() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" {
		port = "8000"
	}

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		printEnvNow(ctx.Writer)
	})
	r.GET("/json", func(ctx *gin.Context) {
		ctx.File("dummy.json")
	})
	r.HEAD("/", func(ctx *gin.Context) {
		printEnvNow(ctx.Writer)
	})
	r.GET("/_stats", gin.WrapF(stats_api.Handler))
	r.Run(host + ":" + port)
}
