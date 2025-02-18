package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Start...")

	store := NewURLStore()

	// regularly check
	go func() {
		for {
			store.CheckURLStatus()
			time.Sleep(30 * time.Second)
		}
	}()

	// router configuration
	http.HandleFunc("/urls", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// get all url status
			statuses := store.GetAll()
			json.NewEncoder(w).Encode(statuses)

		case http.MethodPost:
			// add url
			var request struct{ URL string }
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			id := store.Add(request.URL)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"id": id})

		case http.MethodDelete:
			// remove url
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "Missing ID", http.StatusBadRequest)
				return
			}
			store.Remove(id)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		}
	})

	http.ListenAndServe(":8080", nil)
	fmt.Println("End...")

}
