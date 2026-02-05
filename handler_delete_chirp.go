package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/Bention99/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Bearer Token found", err)
		return
	}

	uID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect bearer Token", err)
		return
	}

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

	if uID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "You have no permission to delete this Chirp.", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithStatus(w, http.StatusNoContent)
}