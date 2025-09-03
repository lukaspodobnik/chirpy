package main

import (
	"encoding/json"
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

func healthzHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add(contenttype, textPlainUTF8)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add(contenttype, textHtml)
	writer.WriteHeader(http.StatusOK)
	writer.Write(fmt.Appendf(
		nil,
		"<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",
		cfg.fileserverHits.Load()))
}

func (cfg *apiConfig) resetHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	writer.WriteHeader(http.StatusOK)
}

func validateChirpHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	var params parameters
	if err := decoder.Decode(&params); err != nil {
		log.Printf("decoding request body failed: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var respBody any
	if len(params.Body) <= 140 {
		respBody = struct {
			Valid bool `json:"valid"`
		}{Valid: true}
		writer.WriteHeader(http.StatusOK)

	} else {
		respBody = struct {
			Error string `json:"error"`
		}{Error: "Chirp is too long"}
		writer.WriteHeader(http.StatusBadRequest)
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("marshalling response body failed: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add(contenttype, appJson)
	writer.Write(data)
}
