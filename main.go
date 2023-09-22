package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	stats_api "github.com/fukata/golang-stats-api-handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Version is to embed the version string
var Version = "202309030000"

//go:embed dummy.json
var embeddedFS embed.FS

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
	var devMode = flag.Bool("dev", true, "enable the proxy for the next dev server")
	flag.Parse()

	if *httpMode {
		HTTP(*devMode)
	} else {
		for range time.Tick(10 * time.Second) {
			printNow(os.Stdout)
		}
	}
}

func RequestIDWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		if reqID != "" {
			w.Header().Add("X-Request-ID", reqID)
		}

		next.ServeHTTP(w, r)
	})
}

// HTTP runs the HTTP server
func HTTP(devMode bool) {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" {
		port = "8000"
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(RequestIDWriter)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			if dur := r.URL.Query().Get("sleep"); dur != "" {
				duration, err := time.ParseDuration(dur)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				time.Sleep(duration)
			}

			w.Header().Add("Content-Type", "application/json")

			printNow(w)
		})

		r.Get("/events", func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				log.Println("INFO: events: client disconnect")
			}()

			flusher, ok := w.(http.Flusher)
			if !ok {
				log.Println("ERROR: no flusher found")
				return
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-store")

			t := time.NewTicker(1 * time.Second)
			defer t.Stop()

			go func() {
				for {
					select {
					case now := <-t.C:
						fmt.Fprintf(w, "data: %d\n\n", now.Unix())
						flusher.Flush()
					}
				}
			}()

			<-r.Context().Done()
		})

		r.Get("/json", staticFileServer(
			http.FS(embeddedFS),
			"dummy.json",
		))

		r.Get("/_stats", stats_api.Handler)
	})

	if devMode {
		log.Print("Setting the proxy for the next dev server...")

		// proxy to the local next dev server
		rpURL, err := url.Parse("http://localhost:3000")
		if err != nil {
			log.Fatal(err)
		}

		r.Handle("/*", &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetXForwarded()
				r.SetURL(rpURL)
			},
		})
	} else {
		r.Handle("/*", http.FileServer(http.Dir("./public")))
	}

	log.Printf("Listening to %s", host+":"+port)

	http.ListenAndServe(host+":"+port, r)
}

func staticFileServer(hfs http.FileSystem, fn string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := hfs.Open(fn)

		switch {
		case errors.Is(err, fs.ErrNotExist):
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		case err != nil:
			http.Error(w, fmt.Errorf("opening file: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			http.Error(w, fmt.Errorf("reading file: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, fn, d.ModTime(), f)
	}
}
