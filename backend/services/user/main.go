package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"database/sql"

	"github.com/ecommerce/observability"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "password123", Role: "customer"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Password: "password123", Role: "admin"},
	{ID: 3, Name: "Demo User", Email: "demo@shop.com", Password: "demo123", Role: "customer"},
}

var dbConn *sql.DB

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"invalid request"}`)
		return
	}

	for _, u := range users {
		if u.Email == req.Email && u.Password == req.Password {
			log.Printf("[AUTH] Login success user=%s email=%s", u.Name, u.Email)
			resp := AuthResponse{
				Token: fmt.Sprintf("token_%d_%d", u.ID, time.Now().Unix()),
				User:  User{ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	log.Printf("[AUTH] Login failed email=%s", req.Email)
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"error":"invalid email or password"}`)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"invalid request"}`)
		return
	}

	for _, u := range users {
		if u.Email == req.Email {
			log.Printf("[AUTH] Register failed - email exists: %s", req.Email)
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, `{"error":"email already registered"}`)
			return
		}
	}

	newUser := User{
		ID:       len(users) + 1,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "customer",
	}

	_, err := dbConn.Exec(`
	INSERT INTO users
	(id,name,email,password,role)
	VALUES($1,$2,$3,$4,$5)
	`,
		newUser.ID,
		newUser.Name,
		newUser.Email,
		newUser.Password,
		newUser.Role,
	)

	if err != nil {
		log.Printf("Insert failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to save user",
		})
		return
	}

	users = append(users, newUser)

	log.Printf("[AUTH] Register success user=%s email=%s", newUser.Name, newUser.Email)
	resp := AuthResponse{
		Token: fmt.Sprintf("token_%d_%d", newUser.ID, time.Now().Unix()),
		User:  User{ID: newUser.ID, Name: newUser.Name, Email: newUser.Email, Role: newUser.Role},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("[USER] Listing %d users", len(users))
	safeUsers := make([]User, len(users))
	for i, u := range users {
		safeUsers[i] = User{ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role}
	}
	json.NewEncoder(w).Encode(safeUsers)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","service":"user-service"}`)
}

func main() {
	mux := http.NewServeMux()
	db := observability.NewTracedDB("user-service")

	dbConn = ConnectDB()

	_, err := dbConn.Exec(`
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL
	)
	`)

	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		_, err := dbConn.Exec(`
			INSERT INTO users
			(id,name,email,password,role)
			VALUES($1,$2,$3,$4,$5)
			ON CONFLICT (id) DO NOTHING
		`,
			u.ID,
			u.Name,
			u.Email,
			u.Password,
			u.Role,
		)

		if err != nil {
			log.Printf("Insert user failed: %v", err)
		}
	}

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/users", getUsers)
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/users/"):]
		db.Query(r.Context(), "users", "SELECT", "id="+id)
		for _, u := range users {
			if fmt.Sprintf("%d", u.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				log.Printf("[USER] Found user id=%s name=%s", id, u.Name)
				json.NewEncoder(w).Encode(User{ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role})
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"user not found"}`)
	})
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		db.Query(r.Context(), "users", "SELECT", "email="+r.URL.Query().Get("email"))
		handleLogin(w, r)
	})
	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r)
		db.Insert(r.Context(), "users", "new_user")
	})

	observability.Run("user-service", mux)
}