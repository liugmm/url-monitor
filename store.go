package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type URLStatus struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	StatusCode  int       `json:"status_code"`
	LastChecked time.Time `json:"last_checked"`
}

type URLStore struct {
	urls map[string]URLStatus
	mu   sync.RWMutex
}

func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]URLStatus),
	}
}

// add url
func (s *URLStore) Add(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.NewString()
	s.urls[id] = URLStatus{
		ID:          id,
		URL:         url,
		StatusCode:  0,
		LastChecked: time.Time{},
	}
	return id
}

// get all status
func (s *URLStore) GetAll() []URLStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	results := make([]URLStatus, 0, len(s.urls))
	for _, v := range s.urls {
		results = append(results, v)
	}
	return results
}
