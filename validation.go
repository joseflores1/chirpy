package main

import (
	"encoding/json"
	"fmt"
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
		respondWithError(w, 500, fmt.Sprintf("couldn't decode request body: %s", errDecode.Error()))
		return
	}

	bodyLen := len(params.Body)
	if bodyLen > 140 {
		respondWithError(w, 400, fmt.Sprintf("Chirp is too long: it has %d characters", bodyLen))
		return
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}
	respondWithJSON(w, 200, returnVals{Valid: true})
}

