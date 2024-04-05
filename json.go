package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	type errorRespone struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errorRespone{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Failed to marshal json response - %v", payload)
		return
	}

	w.Header().Add("Content-Type", "applicaion/json")
	w.WriteHeader(code)
	w.Write(data)
}
