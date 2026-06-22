package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// API Gateway routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"api-gateway"}`)
	})

	// Route ke Product Service
	productProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "product-service:8080",
	})
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		productProxy.ServeHTTP(w, r)
	})

	// Route ke Order Service
	orderProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "order-service:8080",
	})
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		orderProxy.ServeHTTP(w, r)
	})

	// Route ke User Service
	userProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "user-service:8080",
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})

	// Route ke Payment Service
	paymentProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "payment-service:8080",
	})
	mux.HandleFunc("/api/payments", func(w http.ResponseWriter, r *http.Request) {
		paymentProxy.ServeHTTP(w, r)
	})

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
