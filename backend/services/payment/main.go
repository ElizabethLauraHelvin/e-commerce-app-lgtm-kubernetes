package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"created_at"`
}

var payments = []Payment{
	{ID: "1", OrderID: "1", Amount: 999.99, Status: "completed", Method: "credit_card", CreatedAt: time.Now()},
	{ID: "2", OrderID: "2", Amount: 59.98, Status: "pending", Method: "bank_transfer", CreatedAt: time.Now()},
}

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"payment-service"}`)
	})

	// Get all payments
	mux.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payments)
	})

	// Get payment by ID
	mux.HandleFunc("/payments/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/payments/"):]
		for _, p := range payments {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"payment not found"}`)
	})

	// Process payment
	mux.HandleFunc("/payments/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var payment Payment
		if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request"}`)
			return
		}

		payment.ID = fmt.Sprintf("%d", len(payments)+1)
		payment.CreatedAt = time.Now()
		payment.Status = "completed"
		payments = append(payments, payment)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(payment)
	})

	log.Println("Payment Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
