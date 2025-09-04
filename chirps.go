package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lukaspodobnik/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) postChirpHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	var params parameters
	if err := decoder.Decode(&params); err != nil {
		respondWithError(writer, http.StatusInternalServerError, fmt.Sprintf("Decoding request body failed: %v.", err))
		return
	}

	body := params.Body
	uuid, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, fmt.Sprintf("Parsing user_id failed: %v.", err))
		return
	}

	if len(body) <= 140 {
		chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{Body: body, UserID: uuid})
		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, fmt.Sprintf("Creating chirp in db failed: %v.", err))
			return
		}

		respondWithJson(writer, http.StatusCreated, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})

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
