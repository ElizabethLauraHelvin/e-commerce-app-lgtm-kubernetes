package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ecommerce/observability"
)

type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	Method        string    `json:"method"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type PaymentRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
}

var (
	dbConn *sql.DB

	payments []Payment
	paymentCounter = 0
)


func getTransactionPrefix() string {
	if prefix := os.Getenv("TRANSACTION_PREFIX"); prefix != "" {
		return prefix
	}
	return "TXN"
}

func generateTxnID() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return getTransactionPrefix() + "-" + string(b)
}

func main() {

	rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()
	db := observability.NewTracedDB("payment-service")

	dbConn = ConnectDB()

	_, err := dbConn.Exec(`
	CREATE TABLE IF NOT EXISTS payments(
		id TEXT PRIMARY KEY,
		order_id TEXT NOT NULL,
		amount NUMERIC NOT NULL,
		status TEXT NOT NULL,
		method TEXT NOT NULL,
		transaction_id TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	)
	`)

	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","service":"payment-service"}`)
	})

	mux.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			var req PaymentRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Printf("[PAYMENT] Invalid request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"invalid request"}`)
				return
			}

			if req.Method == "" {
				req.Method = "credit_card"
			}

			paymentCounter++
			payment := Payment{
				ID:            fmt.Sprintf("PAY-%05d", paymentCounter),
				OrderID:       req.OrderID,
				Amount:        req.Amount,
				Status:        "completed",
				Method:        req.Method,
				TransactionID: generateTxnID(),
				CreatedAt:     time.Now(),
			}
			payments = append(payments, payment)

			_, err := dbConn.Exec(
				`
				INSERT INTO payments
				(id,order_id,amount,status,method,transaction_id,created_at)
				VALUES($1,$2,$3,$4,$5,$6,$7)
				`,
				payment.ID,
				payment.OrderID,
				payment.Amount,
				payment.Status,
				payment.Method,
				payment.TransactionID,
				payment.CreatedAt,
			)

			if err != nil {
				log.Printf("[PAYMENT] Insert failed: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "failed to save payment",
				})
				return
			}

			log.Printf("[PAYMENT] Processed payment=%s order=%s amount=%.2f method=%s txn=%s",
				payment.ID, payment.OrderID, payment.Amount, payment.Method, payment.TransactionID)
			db.Insert(r.Context(), "payments", payment.ID)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payment)
			return
		}

		// GET - filter by order_id if provided
		orderID := r.URL.Query().Get("order_id")
		if orderID != "" {
			var orderPayments []Payment
			for _, p := range payments {
				if p.OrderID == orderID {
					orderPayments = append(orderPayments, p)
				}
			}
			log.Printf("[PAYMENT] Listing payments for order=%s count=%d", orderID, len(orderPayments))
			json.NewEncoder(w).Encode(orderPayments)
			return
		}

		log.Printf("[PAYMENT] Listing all payments count=%d", len(payments))
		json.NewEncoder(w).Encode(payments)
	})

	observability.Run("payment-service", mux)
}