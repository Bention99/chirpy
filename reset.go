package main

import (
	"net/http"
	"fmt"
	)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.env != "dev" {
		respondWithError(w, http.StatusForbidden, "Can't reset DB outside of dev environment\n", nil)
    	return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset Database\n", err)
    	return
	}
	fmt.Println("Reset sucessfull")
}
