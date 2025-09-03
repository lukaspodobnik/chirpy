package main

import "net/http"

func (cfg *apiConfig) resetHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	writer.WriteHeader(http.StatusOK)
}
