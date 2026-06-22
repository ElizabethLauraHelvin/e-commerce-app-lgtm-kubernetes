package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com", Role: "customer"},
	{ID: "2", Name: "Jane Smith", Email: "jane@example.com", Role: "admin"},
}

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"user-service"}`)
	})

	// Get all users
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	// Get user by ID
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/users/"):]
		for _, u := range users {
			if u.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(u)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"user not found"}`)
	})

	// Create user
	mux.HandleFunc("/users/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request"}`)
			return
		}

		user.ID = fmt.Sprintf("%d", len(users)+1)
		users = append(users, user)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})

	log.Println("User Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
