package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/j-04/url-shortener/storage"
)

func main() {
	mem := storage.NewMemoryStore()
	r := mux.NewRouter()
	redirectPath := "http://localhost:8080/r"
	r.Handle("/add", &AddPath{domain: redirectPath, store: &mem}).Methods("POST")
	r.Handle("/r/{hash}", &DeletePath{store: &mem}).Methods("DELETE")
	r.Handle("/r/{hash}", &RedirectPath{store: &mem}).Methods("GET")
	http.ListenAndServe(":8080", r)
}

type Store interface {
	Add(shortenedUrl, longUrl string) error
	Remove(ShortenedUrl string) error
	Get(shortenedUrl string) (string, error)
}
