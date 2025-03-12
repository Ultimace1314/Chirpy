package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type error struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, error{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			return
		}
		w.Write(data)
	}
}
