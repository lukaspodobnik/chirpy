package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
			CleanedBody string `json:"cleaned_body"`
		}{CleanedBody: removeProfanity(params.Body)})

	} else {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
	}
}

func removeProfanity(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		word = strings.ToLower(word)
		for _, badword := range [3]string{"kerfuffle", "sharbert", "fornax"} {
			if word == badword {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}
