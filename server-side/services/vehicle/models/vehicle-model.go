package models

import (
	"time"
)

type Vehicle struct {
	VehicleID    int    `json:"vehicleID"`
	LicensePlate string `json:"licensePlate"`
	Model        string `json:"model"`
	RentalRate   int    `json:"rentalRate"`
}

type VehicleStatus struct {
	StatusID          int       `json:"statusID"`
	VehicleID         int       `json:"vehicleID"`
	Timestamp         time.Time `json:"timestamp"`
	Location          string    `json:"location"`
	ChargeLevel       int       `json:"chargeLevel"`
	CleanlinessStatus string    `json:"cleanlinessStatus"`
}

type Booking struct {
	BookingID int       `json:"bookingID"`
	VehicleID int       `json:"vehicleID"`
	UserID    int       `json:"userID"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
