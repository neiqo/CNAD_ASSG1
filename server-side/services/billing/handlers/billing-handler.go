package handlers

import (
	"billing/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
func UpdatePaymentStatusToSuccessful(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract paymentID from URL parameters
	paymentID := mux.Vars(r)["paymentID"]

	// Define a struct for the request body to capture userID
	type StatusUpdate struct {
		UserID int `json:"userID"`
	}

	var input StatusUpdate
	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Received input: %+v", input)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Received input: %+v", input)

	// Validate that the userID is provided
	if input.UserID == 0 {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// Update the payment status to 'Successful' for the given paymentID and userID
	query := `UPDATE payment_db.payments SET Status = 'Successful' WHERE paymentID = ? AND userID = ?`
	_, err := db.Exec(query, paymentID, input.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update payment status: %v", err), http.StatusInternalServerError)
		return
	}

	// Retrieve the payment details including the user email, amount, promotion, and booking ID
	var bookingID int
	var userEmail string
	var userName string
	var amount float64
	var promotionName string
	query = `
        SELECT p.bookingID, u.Email, u.Name, p.Amount, pr.Name
        FROM payment_db.payments p
        JOIN users_db.users u ON p.userID = u.userID
        LEFT JOIN common_db.promotions pr ON p.promotionID = pr.promotionID
        WHERE p.paymentID = ?
    `
	err = db.QueryRow(query, paymentID).Scan(&bookingID, &userEmail, &userName, &amount, &promotionName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve payment details: %v", err), http.StatusInternalServerError)
		return
	}

	// Call vehicle-handler to update the booking status to 'Active'
	if err := updateBookingStatusToActive(bookingID, input.UserID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update booking status: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a detailed receipt content
	receiptContent := fmt.Sprintf(
		"Dear %s,\n\n"+
			"Your payment for Booking ID %d has been successfully processed.\n\n"+
			"Payment Details:\n"+
			"- Amount: $%.2f\n"+
			"- Promotion Applied: %s\n"+
			"- Payment Status: Successful\n\n"+
			"Thank you for your payment! Your booking is now active.\n\n"+
			"Best regards,\n"+
			"Vehicle Reservations Team",
		userName, bookingID, amount, promotionName,
	)

	// Send the receipt email to the user
	if err := sendReceiptEmail(userEmail, receiptContent); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send receipt email: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Payment status updated to Successful and booking status updated to Active. Receipt sent to your email."}`))
}

// sendReceiptEmail sends an email with the receipt details to the user
func sendReceiptEmail(toEmail, receiptContent string) error {
	var godotErr = godotenv.Load("../../../.env")
	if godotErr != nil {
		log.Fatalf("Error loading .env file: %v", godotErr)
	}

	GMAIL_EMAIL := os.Getenv("GMAIL_EMAIL")
	GMAIL_PASS := os.Getenv("GMAIL_PASS")

	from := GMAIL_EMAIL
	password := GMAIL_PASS
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Prepare the email message
	subject := "Payment Receipt"
	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, receiptContent))

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Receipt email sent to: %s", toEmail)
	return nil
}

// updateBookingStatusToActive sends a request to the vehicle handler to update the booking status to active
func updateBookingStatusToActive(bookingID, userID int) error {
	// Construct the URL to call the vehicle-handler API
	url := fmt.Sprintf("http://localhost:5002/api/v1/success-payment?bookingID=%d&userID=%d", bookingID, userID)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	// Send the request to the vehicle-handler
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request to vehicle-handler: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to update booking status. Status code: %d", resp.StatusCode)
		return fmt.Errorf("failed to update booking status: %v", resp.Status)
	}

	log.Printf("Booking status updated to Active successfully")
	return nil
}
