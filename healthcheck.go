package main

import (
	"net/http"
	"time"
)

func (s *URLStore) CheckURLStatus() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	concurrency := 10
	sem := make(chan struct{}, concurrency)

	for id, urlStatus := range s.urls {
		sem <- struct{}{}
		go func(id string, url string) {
			defer func() { <-sem }()

			client := http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(url)
			if err != nil {
				// update status to 500 (request failed)
				s.updateStatusCode(id, 500)
				return
			}

			defer resp.Body.Close()

			s.updateStatusCode(id, resp.StatusCode)

		}(id, urlStatus.URL)

	}
}

func (s *URLStore) updateStatusCode(id string, code int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if status, ok := s.urls[id]; ok {
		status.StatusCode = code
		status.LastChecked = time.Now()
		s.urls[id] = status
	}
}
