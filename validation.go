package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
		CleanedBody string `json:"cleaned_body"`
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanBody := getCleanedBody(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanBody,
	})

}

func getCleanedBody(body string, badWords map[string]struct{}) string {

	const replacement = "****"
	isDirty := false

	splitBody := strings.Split(body, " ")

	for i, word := range splitBody {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			splitBody[i] = replacement
			isDirty = true
		}
	}

	if !isDirty {
		return body
	}

	cleanBody := strings.Join(splitBody, " ")
	return cleanBody
}
