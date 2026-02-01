package main

import (
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("ChirpID")
	if chirpIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No Chirp exists by this ID", err)
		return
	}
	respondWithJSON(w, http.StatusOK, answerChirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.CreatedAt.Time,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}
