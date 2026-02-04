package main

import (
	"net/http"
	"github.com/Bention99/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Bearer Token found", err)
		return
	}
	
	err = cfg.db.RevokeToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke Token", err)
		return
	}
	respondWithStatus(w, http.StatusNoContent)
}

func respondWithStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}