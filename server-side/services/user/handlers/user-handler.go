package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"user/models"
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

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body. Error: %v", err), http.StatusBadRequest)
		return
	}

	// check if the email already exists in the database
	var existingEmail string
	emailQuery := "SELECT email FROM users WHERE email = ?"
	err = db.QueryRow(emailQuery, user.Email).Scan(&existingEmail)
	if err == nil {
		// if no error, it means the email already exists
		http.Error(w, "Email is already in use. Please choose a different email.", http.StatusConflict)
		return
	}

	// insert user into the database
	query := "INSERT INTO users (Name, Email, contactNo, hashedPassword) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, user.Name, user.Email, user.ContactNo, user.HashedPassword)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting user into database. Error: %v", err), http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "User registered successfully",
	}
	json.NewEncoder(w).Encode(response)
}
