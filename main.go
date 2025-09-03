package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	filepathroot = "."
	port         = "8080"
)

const (
	contenttype   = "Content-Type"
	textPlainUTF8 = "text/plain; charset=utf-8"
	textHtml      = "text/html"
	appJson       = "application/json"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}

	fileserverHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathroot))))

	mux := http.NewServeMux()

	mux.Handle("/app/", fileserverHandler)
	mux.HandleFunc("GET /api/healthz", healthzHandlerFunc)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandlerFunc)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandlerFunc)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandlerFunc)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %s on port: %s\n", filepathroot, port)
	log.Fatal(server.ListenAndServe())
}
