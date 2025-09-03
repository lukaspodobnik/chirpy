package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) usersHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	var params parameters
	if err := decoder.Decode(&params); err != nil {
		respondWithError(writer, http.StatusInternalServerError, fmt.Sprintf("decoding request body failed: %v", err))
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, fmt.Sprintf("creating user failed: %v", err))
		return
	}

	respondWithJson(writer, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
