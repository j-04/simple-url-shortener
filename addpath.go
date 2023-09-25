package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type addPathRequest struct {
	URL string `json:"url"`
}

type addPathResponse struct {
	ShortenedUrl string `json:"shortened_url"`
	LongUrl      string `json:"long_url"`
}

type AddPath struct {
	domain string
	store  Store
}

func (a *AddPath) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var parsed addPathRequest
	err := json.NewDecoder(r.Body).Decode(&parsed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unexpected error :: %v", err)))
		return
	}

	h := sha1.New()
	h.Write([]byte(parsed.URL))
	sum := h.Sum(nil)
	hash := hex.EncodeToString(sum)[:10]
	err = a.store.Add(hash, parsed.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unexpected error :: %v", err)))
		return
	}

	pathResp := addPathResponse{ShortenedUrl: fmt.Sprintf("%v/%v", a.domain, hash), LongUrl: parsed.URL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pathResp)
}
