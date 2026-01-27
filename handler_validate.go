package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type returnVals struct {
	Valid bool `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	params = wordReplacer(params)

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
		CleanedBody: params.Body,
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