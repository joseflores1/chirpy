package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request body", errDecode)
		return
	}

	const maxChirpLen = 140
	bodyLen := len(params.Body)
	if bodyLen > maxChirpLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil) 
		return
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}
	respondWithJSON(w, http.StatusOK, returnVals{Valid: true})
}

