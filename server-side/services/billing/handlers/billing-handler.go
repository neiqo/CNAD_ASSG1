package handlers

import (
	"billing/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var paymentRequest struct {
		UserID      int     `json:"userID"`
		BookingID   int     `json:"bookingID"`
		Amount      float64 `json:"amount"`
		PromotionID int     `json:"promotionID"`
	}

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid input: %v"}`, err), http.StatusBadRequest)
		return
	}

	log.Printf("promotionid %d", paymentRequest.PromotionID)

	log.Printf("Received payment request: %+v", paymentRequest)
	promotions, err := getPromotionsFromCommonService()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch promotions: %v"}`, err), http.StatusInternalServerError)
		return
	}

	var discount float64
	for _, promo := range promotions {
		if promo.PromotionID == paymentRequest.PromotionID {
			if promo.IfPercentage {
				discount = paymentRequest.Amount * promo.Discount / 100
			} else {
				discount = promo.Discount
			}
			break
		}
	}

	finalAmount := paymentRequest.Amount - discount
	log.Printf("Applied discount of %.2f. Final amount to pay: %.2f", discount, finalAmount)

	query := `INSERT INTO payment_db.payments (userID, bookingID, Amount, Status, promotionID) 
              VALUES (?, ?, ?, 'Pending', ?)`
	result, err := db.Exec(query, paymentRequest.UserID, paymentRequest.BookingID, finalAmount, paymentRequest.PromotionID)

	if err != nil {
		log.Printf("Error executing SQL query: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create payment record. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted ID: %v", err)
		http.Error(w, `{"error": "Failed to retrieve last inserted ID"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Payment created successfully with ID: %d", lastInsertID)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"message": "Payment created and is pending!", "paymentID": %d, "finalAmount": %.2f}`, lastInsertID, finalAmount)))
}

// getPromotionsFromCommonService makes a request to the Common Service to fetch all promotions.
func getPromotionsFromCommonService() ([]models.Promotion, error) {
	resp, err := http.Get("http://localhost:5003/api/v1/promotions")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch promotions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch promotions. Status code: %d", resp.StatusCode)
	}

	var promotions []models.Promotion
	if err := json.NewDecoder(resp.Body).Decode(&promotions); err != nil {
		return nil, fmt.Errorf("failed to decode promotions: %v", err)
	}

	return promotions, nil
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
