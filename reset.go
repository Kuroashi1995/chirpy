package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not Authorized")
	}
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	err := cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users")
		return
	}
	w.Write([]byte("Hits reset to 0, users deleted"))
}
