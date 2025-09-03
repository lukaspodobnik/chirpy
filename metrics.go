package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricsHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add(contenttype, textHtml)
	writer.WriteHeader(http.StatusOK)
	writer.Write(fmt.Appendf(
		nil,
		"<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",
		cfg.fileserverHits.Load()))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
