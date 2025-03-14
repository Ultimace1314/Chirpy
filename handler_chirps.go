package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ultimace1314/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}

func (cfg apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cfg.handlerGetChirps(w, r)
	case http.MethodPost:
		cfg.handlerPostChirps(w, r)
	default:
		respondWithError(w, 405, "Method not allowed", nil)
	}
}

func (cfg apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Couldn't get chirps", err)
		return
	}

	chirplist := []Chirp{}

	for _, chirp := range chirps {
		c := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			User_id:   chirp.UserID,
		}
		chirplist = append(chirplist, c)
	}

	respondWithJSON(w, http.StatusOK, chirplist)
}

func (cfg apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := Params{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}
	if len(params.Body) == 0 {
		respondWithError(w, 400, "Chirp is empty", err)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	post, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.User_id,
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      post.Body,
		User_id:   post.UserID,
	})
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg apiConfig) handlerChirpsByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), uuidID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	})
}
