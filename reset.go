package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if platform := os.Getenv("PLATFORM"); platform != "dev" {
		http.Error(w, "Reset not allowed outside of dev environment.", http.StatusForbidden)
		return
	} else {
		cfg.db.Reset(r.Context())
		w.Write([]byte("Users table has been cleared\n"))
		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0\n"))
	}
}
