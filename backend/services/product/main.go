package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price float64 `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: "1", Name: "Laptop", Price: 999.99, Stock: 5},
	{ID: "2", Name: "Mouse", Price: 29.99, Stock: 50},
	{ID: "3", Name: "Keyboard", Price: 79.99, Stock: 30},
}

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"product-service"}`)
	})

	// Get all products
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	// Get product by ID
	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/products/"):]
		for _, p := range products {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"product not found"}`)
	})

	log.Println("Product Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
