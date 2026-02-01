package main

import (
	"net/http"
	"time"
	"github.com/google/uuid"
)

type answerChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     	string 	`json:"body"`
	UserID	uuid.UUID	`json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	allChirps := []answerChirp{}

	for _, chirp := range chirps {
		allChirps = append(allChirps, answerChirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.CreatedAt.Time,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}
