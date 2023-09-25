package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type RedirectPath struct {
	store Store
}

func (p *RedirectPath) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("empty hash"))
		return
	}

	longUrl, err := p.store.Get(hash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	http.Redirect(w, r, longUrl, http.StatusTemporaryRedirect)
}

type internalStore struct {
	Version string            `json:"version"`
	Items   map[string]string `json:"items"`
}

type FileStore struct {
	filename string
}

func NewFileStore(filename string) (FileStore, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		is := internalStore{Version: "v1", Items: make(map[string]string)}
		raw, err := json.Marshal(is)
		if err != nil {
			return FileStore{}, fmt.Errorf("unable to generate json representation for file")
		}
		err = os.WriteFile(filename, raw, 0644)
		if err != nil {
			return FileStore{}, fmt.Errorf("Unable to persist file")
		}
	}
	return FileStore{filename: filename}, nil
}
