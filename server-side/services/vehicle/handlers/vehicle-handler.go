package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"vehicle/models"
)

var db *sql.DB

func SetDBConnection(database *sql.DB) {
	db = database
}

func Status(getDBStatus func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !getDBStatus() {
			http.Error(w, "Error: Vehicle Service failed to connect to the database", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Vehicle Service connected to the database successfully!")
	}
}

func AddVehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type VehicleInput struct {
		LicensePlate string `json:"licensePlate"`
		Model        string `json:"model"`
		RentalRate   int    `json:"rentalRate"`
	}

	var input VehicleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO vehicles (licensePlate, Model, rentalRate) VALUES (?, ?, ?)`
	result, err := db.Exec(query, input.LicensePlate, input.Model, input.RentalRate)
	if err != nil {
		http.Error(w, "Failed to insert vehicle", http.StatusInternalServerError)
		return
	}

	insertID, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Vehicle added successfully with ID: %d", insertID)
}

func AddVehicleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type VehicleStatusInput struct {
		VehicleID         int    `json:"vehicleID"`
		Location          string `json:"location"`
		ChargeLevel       int    `json:"chargeLevel"`
		CleanlinessStatus string `json:"cleanlinessStatus"`
	}

	var input VehicleStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO vehicleStatusHistory (vehicleID, location, chargeLevel, cleanlinessStatus) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, input.VehicleID, input.Location, input.ChargeLevel, input.CleanlinessStatus)
	if err != nil {
		http.Error(w, "Failed to insert vehicle status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Vehicle status added successfully!")
}

func AddBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	type BookingInput struct {
		VehicleID int       `json:"vehicleID"`
		UserID    int       `json:"userID"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	}

	var input BookingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid input: %v"}`, err), http.StatusBadRequest)
		log.Printf("Error decoding input: %v", err)
		return
	}

	query := `
        SELECT COUNT(*) 
        FROM bookings 
        WHERE vehicleID = ? 
        AND (
            (startTime BETWEEN ? AND ?) OR 
            (endTime BETWEEN ? AND ?) OR 
            (? BETWEEN startTime AND endTime) OR 
            (? BETWEEN startTime AND endTime)
        )`
	var count int
	err := db.QueryRow(query, input.VehicleID, input.StartTime, input.EndTime, input.StartTime, input.EndTime, input.StartTime, input.EndTime).Scan(&count)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error checking availability for the requested time slot: %v"}`, err), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, `{"error": "Vehicle is already booked for the requested time slot"}`, http.StatusBadRequest)
		return
	}

	query = `INSERT INTO bookings (vehicleID, userID, startTime, endTime) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, input.VehicleID, input.UserID, input.StartTime, input.EndTime)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create booking. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "Booking created successfully!"}`)
}
func GetVehicles(w http.ResponseWriter, r *http.Request) {
	query := `SELECT vehicleID, licensePlate, Model, rentalRate FROM vehicles`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []map[string]interface{}
	for rows.Next() {
		var id int
		var plate, model string
		var rate int
		rows.Scan(&id, &plate, &model, &rate)
		vehicles = append(vehicles, map[string]interface{}{
			"vehicleID":    id,
			"licensePlate": plate,
			"model":        model,
			"rentalRate":   rate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func GetVehicleByID(w http.ResponseWriter, r *http.Request) {
	vehicleID := r.URL.Query().Get("vehicleID")
	if vehicleID == "" {
		http.Error(w, "vehicleID is required", http.StatusBadRequest)
		return
	}

	var vehicle models.Vehicle
	var status models.VehicleStatus
	var timestampStr string

	queryVehicle := `
        SELECT vehicleID, licensePlate, Model, rentalRate
        FROM vehicles
        WHERE vehicleID = ?`
	err := db.QueryRow(queryVehicle, vehicleID).Scan(
		&vehicle.VehicleID,
		&vehicle.LicensePlate,
		&vehicle.Model,
		&vehicle.RentalRate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch vehicle details. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	queryStatus := `
        SELECT statusID, vehicleID, timestamp, location, chargeLevel, cleanlinessStatus
        FROM vehicleStatusHistory
        WHERE vehicleID = ?
        ORDER BY timestamp DESC
        LIMIT 1`
	err = db.QueryRow(queryStatus, vehicleID).Scan(
		&status.StatusID,
		&status.VehicleID,
		&timestampStr,
		&status.Location,
		&status.ChargeLevel,
		&status.CleanlinessStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]interface{}{
				"vehicle": vehicle,
				"status":  nil,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch vehicle status. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	status.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to parse timestamp. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"vehicle": vehicle,
		"status":  status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func GetPastBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, `{"error": "userID is required"}`, http.StatusBadRequest)
		return
	}

	query := `
        SELECT bookingID, vehicleID, startTime, endTime, status
        FROM bookings
        WHERE userID = ? AND status IN ('Completed', 'Cancelled')
        ORDER BY endTime DESC
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error fetching past bookings: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []map[string]interface{}
	for rows.Next() {
		var bookingID, vehicleID int
		var startTimeStr, endTimeStr, status string

		err := rows.Scan(&bookingID, &vehicleID, &startTimeStr, &endTimeStr, &status)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error scanning booking: %v"}`, err), http.StatusInternalServerError)
			return
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error parsing startTime: %v"}`, err), http.StatusInternalServerError)
			return
		}
		endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error parsing endTime: %v"}`, err), http.StatusInternalServerError)
			return
		}

		bookings = append(bookings, map[string]interface{}{
			"bookingID": bookingID,
			"vehicleID": vehicleID,
			"startTime": startTime,
			"endTime":   endTime,
			"status":    status,
		})
	}

	if len(bookings) == 0 {
		http.Error(w, `{"error": "No completed or cancelled bookings found for the specified user"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
func GetUpcomingBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, `{"error": "userID is required"}`, http.StatusBadRequest)
		return
	}

	// Get the current time
	currentTime := time.Now()

	query := `
        SELECT 
            b.bookingID, 
            v.vehicleID, 
            v.licensePlate, 
            v.Model, 
            v.rentalRate, 
            b.startTime, 
            b.endTime, 
            b.status
        FROM bookings b
        JOIN vehicles v ON b.vehicleID = v.vehicleID
        WHERE b.userID = ? AND b.startTime > ? AND b.status = 'Active'
        ORDER BY b.startTime ASC
    `

	rows, err := db.Query(query, userID, currentTime)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error fetching upcoming bookings: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []map[string]interface{}
	for rows.Next() {
		var bookingID, vehicleID, rentalRate int
		var licensePlate, model, status string
		var startTimeStr, endTimeStr string

		// Scan each row into variables
		err := rows.Scan(&bookingID, &vehicleID, &licensePlate, &model, &rentalRate, &startTimeStr, &endTimeStr, &status)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error scanning booking: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Parse the startTime and endTime from strings to time.Time
		startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error parsing startTime: %v"}`, err), http.StatusInternalServerError)
			return
		}
		endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error parsing endTime: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Add booking and vehicle details to the response
		bookings = append(bookings, map[string]interface{}{
			"bookingID":    bookingID,
			"vehicleID":    vehicleID,
			"licensePlate": licensePlate,
			"model":        model,
			"rentalRate":   rentalRate,
			"startTime":    startTime,
			"endTime":      endTime,
			"status":       status,
		})
	}

	// Check if there are any upcoming bookings
	if len(bookings) == 0 {
		http.Error(w, `{"error": "No upcoming bookings found for the specified user"}`, http.StatusNotFound)
		return
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func CancelBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	bookingID := r.URL.Query().Get("bookingID")
	userID := r.URL.Query().Get("userID")
	if bookingID == "" || userID == "" {
		http.Error(w, `{"error": "bookingID and userID are required"}`, http.StatusBadRequest)
		return
	}

	var startTimeStr string
	var status string
	query := `
        SELECT startTime, status 
        FROM bookings 
        WHERE bookingID = ? AND userID = ?`
	err := db.QueryRow(query, bookingID, userID).Scan(&startTimeStr, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Booking not found or does not belong to the user"}`, http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Error checking booking details: %v"}`, err), http.StatusInternalServerError)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error parsing startTime: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if startTime.Before(time.Now()) {
		http.Error(w, `{"error": "Cannot cancel a booking that has already started or completed"}`, http.StatusBadRequest)
		return
	}

	updateQuery := `
        UPDATE bookings 
        SET status = 'Cancelled' 
        WHERE bookingID = ? AND userID = ? AND status != 'Cancelled'`
	result, err := db.Exec(updateQuery, bookingID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error updating booking status: %v"}`, err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error checking affected rows: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "Booking was already cancelled or doesn't belong to the user"}`, http.StatusBadRequest)
		return
	}

	// Return success message as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Booking successfully cancelled",
	})
}
