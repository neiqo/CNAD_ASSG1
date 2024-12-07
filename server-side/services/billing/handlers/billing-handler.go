package handlers

import (
	"billing/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var db *sql.DB

func SetDBConnection(database *sql.DB) {
	db = database
}

func Status(getDBStatus func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !getDBStatus() {
			http.Error(w, "Error: Billing Service failed to connect to the database", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Billing Service connected to the database successfully!")
	}
}

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type PaymentInput struct {
		UserID      int     `json:"userID"`
		BookingID   int     `json:"bookingID"`
		PromotionID int     `json:"promotionID"`
		Amount      float64 `json:"amount"`
	}

	var input PaymentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO payments (userID, bookingID, promotionID, Amount) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, input.UserID, input.BookingID, input.PromotionID, input.Amount)
	if err != nil {
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	paymentID, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"paymentID": ` + string(paymentID) + `}`))
}

func UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type StatusUpdate struct {
		Status string `json:"status"`
	}

	paymentID := mux.Vars(r)["paymentID"]
	var input StatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE payments SET Status = ? WHERE paymentID = ?`
	_, err := db.Exec(query, input.Status, paymentID)
	if err != nil {
		http.Error(w, "Failed to update payment status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Payment status updated successfully"}`))
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	paymentID := mux.Vars(r)["paymentID"]

	var payment models.Payment
	var createdAtStr string

	query := `SELECT paymentID, userID, bookingID, Status, promotionID, Amount, createdAt FROM payments WHERE paymentID = ?`
	err := db.QueryRow(query, paymentID).Scan(
		&payment.PaymentID,
		&payment.UserID,
		&payment.BookingID,
		&payment.Status,
		&payment.PromotionID,
		&payment.Amount,
		&createdAtStr,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve payment details. Error: %v", err), http.StatusInternalServerError)
		return
	}

	payment.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse createdAt field. Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}
