package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	if len(params.Body) == 0 {
		respondWithError(w, 400, "Chirp is empty")
		return
	}

	retResp := handlerBadWordCleaner(params.Body)

	respondWithJSON(w, 200, response{CleanedBody: retResp})
}

func handlerBadWordCleaner(s string) string {
	badwords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(s, " ")
	// for each word in the string
	for i, word := range words {
		// for each bad word
		for _, badword := range badwords {
			// if the word is a bad word
			if strings.ToLower(word) == badword {
				// replace the word with ****
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
