package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

	if len(params.Body) <= 140 {
		respondWithJson(writer, http.StatusOK, struct {
			Valid bool `json:"valid"`
		}{Valid: true})

	} else {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
	}
}
