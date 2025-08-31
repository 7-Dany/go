package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errResoponse struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal json response:%v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "Application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	respondWithJSON(w, code, errResoponse{Error: msg})
}
