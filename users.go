package main

import (
	"encoding/json"
	"log"
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
		log.Printf("decoding request body failed: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	userDB, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("creating user failed: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := User{
		ID:        userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email:     userDB.Email,
	}

	respondWithJson(writer, http.StatusCreated, user)
}
