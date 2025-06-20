package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", errDecode)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, errUpdate := cfg.db.UpgradeMembership(r.Context(), params.Data.UserID)
	if errUpdate != nil {
		if errUpdate == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "User doesn't exist", errUpdate)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", errUpdate)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
