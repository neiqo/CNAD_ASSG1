package main

import (
	"database/sql"
	"vehicle/handlers"

	//"vehicle/models
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var db *sql.DB
var dbConnected int32

func main() {
	// DB CONNECTION
	var err error

	err = godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("VehicleServiceErr: Error loading .env file: %v", err)
	}

	DB_AUTH := os.Getenv("DB_AUTH")

	db, err = sql.Open("mysql", DB_AUTH+"vehicles_reservations_db")
	if err != nil {
		log.Printf("Initial database connection failed: %v", err)
	}

	// start a goroutine to retry the connection if it fails
	go retryDBConnection(DB_AUTH)

	handlers.SetDBConnection(db)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// ROUTES
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/status", handlers.Status(getDBStatus)) // Fallback Status Route

	router.HandleFunc("/api/v1/vehicles", handlers.AddVehicle).Methods("POST")
	router.HandleFunc("/api/v1/vehicles/status", handlers.AddVehicleStatus).Methods("POST")
	router.HandleFunc("/api/v1/bookings", handlers.AddBooking).Methods("POST")
	router.HandleFunc("/api/v1/vehicles", handlers.GetVehicles).Methods("GET")
	router.HandleFunc("/api/v1/vehicle", handlers.GetVehicleByID).Methods("GET")
	router.HandleFunc("/api/v1/past-bookings", handlers.GetPastBookings).Methods("GET")
	router.HandleFunc("/api/v1/upcoming-bookings", handlers.GetUpcomingBookings).Methods("GET")
	router.HandleFunc("/api/v1/cancel-booking", handlers.CancelBooking).Methods("PUT")
	router.HandleFunc("/api/v1/modify-booking", handlers.ModifyBooking).Methods("PUT")
	router.HandleFunc("/api/v1/success-payment", handlers.UpdateBookingStatus).Methods("PUT")

	fmt.Println("Vehicle Service listening at port 5002")
	log.Fatal(http.ListenAndServe("localhost:5002", c.Handler(router)))

}

func retryDBConnection(dbAuth string) {
	for {
		if err := db.Ping(); err != nil {
			log.Printf("Vehicle Database connection failed: %v. Retrying in 5 seconds...", err)

			// keep trying every 5 seconds reconnect to the db
			var err error
			db, err = sql.Open("mysql", dbAuth+"vehicles_reservations_db")
			if err != nil {
				log.Printf("Error reinitializing database connection: %v", err)
			}
		} else {
			log.Println("Vehicle Database connection successful!")
			SetDBConnected(true)
			return
		}

		// wait 5 seconds before retrying
		time.Sleep(5 * time.Second)
	}
}

func SetDBConnected(status bool) {
	if status {
		atomic.StoreInt32(&dbConnected, 1)
	} else {
		atomic.StoreInt32(&dbConnected, 0)
	}
}

func getDBStatus() bool {
	return atomic.LoadInt32(&dbConnected) == 1
}
