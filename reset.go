package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) postResetHandler(writer http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(writer, http.StatusForbidden, "Reset only allowed in dev environment.")
		return
	}

	if err := cfg.dbQueries.DeleteAllUsers(req.Context()); err != nil {
		respondWithError(writer, http.StatusInternalServerError, fmt.Sprintf("Resetting database failed: %v.", err))
		return
	}

	cfg.fileserverHits.Store(0)
	respondWithJson(writer, http.StatusOK, "Hits reset to 0 and database reset to inital state.")
}
