package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	stats_api "github.com/fukata/golang-stats-api-handler"
	"github.com/gin-gonic/gin"
)

// Version is to embed the version string
var Version = "201601170013"

type response struct {
	Version string    `json:"version"`
	Now     time.Time `json:"now"`
}

func printNow(w io.Writer) {
	err := json.NewEncoder(w).Encode(&response{
		Version: Version,
		Now:     time.Now(),
	})

	if err != nil {
		log.Print("ERROR:", err)
	}
}

func main() {
	var httpMode = flag.Bool("http", true, "enable HTTP server")
	flag.Parse()

	if *httpMode {
		HTTP()
	} else {
		for range time.Tick(10 * time.Second) {
			printNow(os.Stdout)
		}
	}
}

// HTTP runs the HTTP server
func HTTP() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" {
		port = "8000"
	}

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		if dur := ctx.Query("sleep"); dur != "" {
			duration, err := time.ParseDuration(dur)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
			time.Sleep(duration)
		}
		printNow(ctx.Writer)
	})
	r.GET("/json", func(ctx *gin.Context) {
		ctx.File("dummy.json")
	})
	r.HEAD("/json", func(ctx *gin.Context) {
		ctx.File("dummy.json")
	})
	r.HEAD("/", func(ctx *gin.Context) {
		printNow(ctx.Writer)
	})
	r.GET("/_stats", gin.WrapF(stats_api.Handler))
	r.Run(host + ":" + port)
}
