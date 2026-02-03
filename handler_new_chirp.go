package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/Bention99/chirpy/internal/database"
	"github.com/Bention99/chirpy/internal/auth"
)

type parameters struct {
	Body string `json:"body"`
	UID uuid.UUID `json:"user_id"`
}

type returnVals struct {
	Valid bool `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

type createdChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     	string 	`json:"body"`
	UserID	uuid.UUID	`json:"user_id"`
}

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No bearer Token found", err)
		return
	}

	uID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect bearer Token", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	params = wordReplacer(params)

	chirpParams := database.CreateChirpParams{
		Body: params.Body,
		UserID: uID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdChirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.CreatedAt.Time,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func wordReplacer(p parameters) parameters {
	words := strings.Split(p.Body, " ")
	cleanedWords := []string{}
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range words {
		cleaned := strings.ToLower(word)
		bad := false
		for _, badWord := range badWords {
			if cleaned == badWord {
			cleaned = "****"
			bad = true
			}
		}
		if bad {
			cleanedWords = append(cleanedWords, cleaned)
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	newBody := strings.Join(cleanedWords, " ")
	p.Body = newBody
	return p
}