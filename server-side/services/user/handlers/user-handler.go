package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
	"user/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func SetDBConnection(database *sql.DB) {
	db = database
}

func Status(getDBStatus func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !getDBStatus() {
			http.Error(w, "Error: User Service failed to connect to the database", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User Service connected to the database successfully!")
	}
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func sendVerificationEmail(toEmail, verificationCode string) error {
	var godotErr = godotenv.Load("../../../.env")
	if godotErr != nil {
		log.Fatalf("UserServiceErr: Error loading .env file: %v", godotErr)
	}

	GMAIL_EMAIL := os.Getenv("GMAIL_EMAIL")
	GMAIL_PASS := os.Getenv("GMAIL_PASS")

	from := GMAIL_EMAIL
	password := GMAIL_PASS
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte(fmt.Sprintf("Subject: Email Verification\n\nYour verification code is: %s", verificationCode))
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request body. Error: %v"}`, err), http.StatusBadRequest)
		return
	}

	// check if the email already exists in the database
	var existingEmail string
	emailQuery := "SELECT email FROM users WHERE email = ?"
	err = db.QueryRow(emailQuery, user.Email).Scan(&existingEmail)
	if err == nil {
		// if no error, it means the email already exists
		http.Error(w, `{"error": "Email is already in use. Please choose a different email."}`, http.StatusConflict)
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error hashing password. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	verificationCode := generateVerificationCode()

	// send the verification email
	err = sendVerificationEmail(user.Email, verificationCode)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error sending verification email. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// save the user with hashed password and verification code in the database
	query := "INSERT INTO users (Name, Email, contactNo, hashedPassword, verifCode, emailVerified) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, user.Name, user.Email, user.ContactNo, string(hashedPassword), verificationCode, false)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error inserting user into database. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "User registered successfully. A verification code has been sent to your email.",
	}
	json.NewEncoder(w).Encode(response)
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verificationCode"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request body. Error: %v"}`, err), http.StatusBadRequest)
		return
	}

	var storedCode string
	var emailVerified bool

	// query the database to get the stored verification code and email verified status
	query := "SELECT verifCode, emailVerified FROM users WHERE email = ?"
	err = db.QueryRow(query, request.Email).Scan(&storedCode, &emailVerified)
	if err != nil {
		// if no user is found
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
			return
		}
		// if database error
		http.Error(w, fmt.Sprintf(`{"error": "Error querying user from database. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// if the email already verified
	if emailVerified {
		http.Error(w, `{"error": "Email already verified"}`, http.StatusBadRequest)
		return
	}

	// compare the provided verification code with the stored code
	if storedCode != request.VerificationCode {
		// if the verification code doesn't match
		http.Error(w, `{"error": "Invalid verification code"}`, http.StatusUnauthorized)
		return
	}

	// if the verification code is valid, update the emailVerified status in the database
	updateQuery := "UPDATE users SET emailVerified = true WHERE email = ?"
	_, err = db.Exec(updateQuery, request.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error updating email verification status. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Email successfully verified",
	}
	json.NewEncoder(w).Encode(response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email          string `json:"email"`
		HashedPassword string `json:"hashedPassword"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request body. Error: %v"}`, err), http.StatusBadRequest)
		return
	}
	log.Printf("Inhgeathaeth:%v", loginData.Email)
	log.Printf("Inhgeathaeth:%v", loginData.HashedPassword)

	var user models.User

	// query user from the database
	query := "SELECT userID, Name, Email, contactNo, hashedPassword, membershipTier, emailVerified FROM users WHERE email = ?"
	err = db.QueryRow(query, loginData.Email).Scan(&user.UserID, &user.Name, &user.Email, &user.ContactNo, &user.HashedPassword, &user.MembershipTier, &user.EmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Error querying user from database. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// check if email is verified
	if !user.EmailVerified {
		http.Error(w, `{"error": "Email not verified. Please verify your email before logging in."}`, http.StatusUnauthorized)
		return
	}

	// compare the password received with the stored bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginData.HashedPassword))
	if err != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Login successful",
	}
	json.NewEncoder(w).Encode(response)
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	email = strings.TrimSpace(email)

	// check if email is provided
	if email == "" {
		http.Error(w, `{"error": "Email parameter is required"}`, http.StatusBadRequest)
		return
	}

	var user models.User

	// query the database to fetch user details by email
	query := "SELECT userID, Name, Email, contactNo, hashedPassword, membershipTier, emailVerified FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.UserID, &user.Name, &user.Email, &user.ContactNo, &user.HashedPassword, &user.MembershipTier, &user.EmailVerified)

	if err != nil {
		// if no user is found
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
			return
		}
		// if there is an error querying the database
		http.Error(w, fmt.Sprintf(`{"error": "Error querying user from database. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// return user details in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	email := vars["email"] // Extract user email from URL params

	email = strings.TrimSpace(email)

	// Parse the JSON request body
	var updateData struct {
		Name      string `json:"name,omitempty"`
		ContactNo string `json:"contactNo,omitempty"`
		Password  string `json:"password,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request body. Error: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Prepare to update the fields that were provided in the request
	var updateFields []string
	var updateValues []interface{}

	// Check and add name if it's provided
	if updateData.Name != "" {
		updateFields = append(updateFields, "Name = ?")
		updateValues = append(updateValues, updateData.Name)
	}

	// Check and add contact number if it's provided
	if updateData.ContactNo != "" {
		updateFields = append(updateFields, "contactNo = ?")
		updateValues = append(updateValues, updateData.ContactNo)
	}

	// Check and add password if it's provided
	if updateData.Password != "" {
		// Hash the new password before storing it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error hashing new password. Error: %v"}`, err), http.StatusInternalServerError)
			return
		}
		updateFields = append(updateFields, "hashedPassword = ?")
		updateValues = append(updateValues, string(hashedPassword))
	}

	// If no fields were provided, return an error
	if len(updateFields) == 0 {
		http.Error(w, `{"error": "No fields provided to update"}`, http.StatusBadRequest)
		return
	}

	// Construct the update query
	query := fmt.Sprintf("UPDATE users SET %s WHERE email = ?", strings.Join(updateFields, ", "))
	updateValues = append(updateValues, email)

	// Execute the query
	_, err = db.Exec(query, updateValues...)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error updating user details. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "User details updated successfully.",
	}
	json.NewEncoder(w).Encode(response)
}
