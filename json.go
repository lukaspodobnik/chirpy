package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	respondWithJson(writer, code, struct {
		Error string `json:"error"`
	}{Error: msg})
}

func respondWithJson(writer http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("marshalling failed: %v\n", err)
		return
	}

	writer.Header().Set(contenttype, appJson)
	writer.WriteHeader(code)
	if _, err := writer.Write(data); err != nil {
		log.Printf("writing response body failed: %v\n", err)
	}
}
