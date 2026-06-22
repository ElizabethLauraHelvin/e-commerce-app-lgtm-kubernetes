package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ecommerce/observability"
)

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Category string  `json:"category"`
	Image    string  `json:"image"`
}

var products = []Product{
	{ID: "1", Name: "MacBook Pro 14\" M3", Price: 1999.99, Stock: 15, Category: "Laptop", Image: "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400&h=400&fit=crop"},
	{ID: "2", Name: "iPhone 15 Pro Max", Price: 1199.99, Stock: 30, Category: "Smartphone", Image: "https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400&h=400&fit=crop"},
	{ID: "3", Name: "Samsung Galaxy S24 Ultra", Price: 1099.99, Stock: 25, Category: "Smartphone", Image: "https://images.unsplash.com/photo-1610945265064-0e34e5519bbf?w=400&h=400&fit=crop"},
	{ID: "4", Name: "Sony WH-1000XM5", Price: 349.99, Stock: 40, Category: "Audio", Image: "https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?w=400&h=400&fit=crop"},
	{ID: "5", Name: "iPad Air M2", Price: 799.99, Stock: 20, Category: "Tablet", Image: "https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=400&h=400&fit=crop"},
	{ID: "6", Name: "Logitech MX Master 3S", Price: 99.99, Stock: 60, Category: "Accessories", Image: "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=400&h=400&fit=crop"},
	{ID: "7", Name: "Samsung 4K Smart TV 55\"", Price: 699.99, Stock: 10, Category: "Electronics", Image: "https://images.unsplash.com/photo-1593359677879-a4bb92f829d1?w=400&h=400&fit=crop"},
	{ID: "8", Name: "Nike Air Max 270", Price: 149.99, Stock: 45, Category: "Fashion", Image: "https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=400&h=400&fit=crop"},
	{ID: "9", Name: "Mechanical Keyboard RGB", Price: 129.99, Stock: 35, Category: "Accessories", Image: "https://images.unsplash.com/photo-1618384887929-16ec33fab9ef?w=400&h=400&fit=crop"},
	{ID: "10", Name: "GoPro Hero 12", Price: 399.99, Stock: 18, Category: "Electronics", Image: "https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=400&h=400&fit=crop"},
	{ID: "11", Name: "Adidas Ultraboost 23", Price: 189.99, Stock: 50, Category: "Fashion", Image: "https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=400&h=400&fit=crop"},
	{ID: "12", Name: "Nintendo Switch OLED", Price: 349.99, Stock: 22, Category: "Gaming", Image: "https://images.unsplash.com/photo-1606144042614-b2417e99c4e3?w=400&h=400&fit=crop"},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"product-service"}`)
	})

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Printf("[PRODUCT] Returning %d products", len(products))
		json.NewEncoder(w).Encode(products)
	})

	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/products/"):]
		for _, p := range products {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				log.Printf("[PRODUCT] Found product id=%s name=%s", p.ID, p.Name)
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		log.Printf("[PRODUCT] Product not found id=%s", id)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"product not found"}`)
	})

	observability.Run("product-service", mux)
}
