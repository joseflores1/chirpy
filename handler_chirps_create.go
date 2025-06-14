package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	User_ID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpCreation(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", errDecode)
		return
	}

	cleanBody, errValidate := validateChirp(params.Body)
	if errValidate != nil {
		respondWithError(w, http.StatusBadRequest, errValidate.Error(), errValidate)
		return
	}

	chirp, errCreateChirp := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanBody,
		UserID: params.User_ID,
	})
	if errCreateChirp != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", errCreateChirp)
		return
	}

	type response struct {
		Chirp
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			User_ID: chirp.UserID,
		},
	})
}

func validateChirp(body string) (string, error) {

	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
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
