package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var orders = []Order{
	{ID: "1", UserID: "user1", ProductID: "1", Quantity: 1, Total: 999.99, Status: "completed", CreatedAt: time.Now()},
	{ID: "2", UserID: "user2", ProductID: "2", Quantity: 2, Total: 59.98, Status: "pending", CreatedAt: time.Now()},
}

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"order-service"}`)
	})

	// Get all orders
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})

	// Get order by ID
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/orders/"):]
		for _, o := range orders {
			if o.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(o)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"order not found"}`)
	})

	// Create order
	mux.HandleFunc("/orders/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request"}`)
			return
		}

		order.ID = fmt.Sprintf("%d", len(orders)+1)
		order.CreatedAt = time.Now()
		orders = append(orders, order)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	})

	log.Println("Order Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
