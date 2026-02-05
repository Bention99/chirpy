package main

import (
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/Bention99/chirpy/internal/database"
)

type answerChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     	string 	`json:"body"`
	UserID	uuid.UUID	`json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")

	var (
		chirps []database.Chirp
		err    error
	)

	if authorIDStr != "" {
		authorUUID, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", parseErr)
			return
		}

		chirps, err = cfg.db.GetChirpsFromUser(r.Context(), authorUUID)
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	allChirps := make([]answerChirp, 0, len(chirps))
	for _, chirp := range chirps {
		allChirps = append(allChirps, answerChirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}
