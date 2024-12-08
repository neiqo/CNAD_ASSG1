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

	log.Printf("VehicleID: %d", input.VehicleID)
	log.Printf("UserID: %d", input.UserID)
	log.Printf("Start Time: %v", input.StartTime)
	log.Printf("End Time: %v", input.EndTime)

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
