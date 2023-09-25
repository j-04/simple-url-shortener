package storage

import (
	"fmt"
	"log"
)

type MemoryStore struct {
	items map[string]string
}

func NewMemoryStore() MemoryStore {
	return MemoryStore{items: make(map[string]string)}
}

func (m *MemoryStore) Add(shortenedUrl, longUrl string) error {
	if m.items[shortenedUrl] != "" {
		return fmt.Errorf("Value already exists here")
	}
	m.items[shortenedUrl] = longUrl
	log.Println(m.items)
	return nil
}

func (m *MemoryStore) Remove(shortenedUrl string) error {
	if m.items[shortenedUrl] == "" {
		return fmt.Errorf("Value does not exist here")
	}
	delete(m.items, shortenedUrl)
	return nil
}

func (m *MemoryStore) Get(shortenedUrl string) (string, error) {
	longUrl, ok := m.items[shortenedUrl]
	if !ok {
		return "", fmt.Errorf("No mapped url available here")
	}
	return longUrl, nil
}
