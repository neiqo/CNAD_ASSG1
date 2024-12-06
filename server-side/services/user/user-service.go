package main

import (
	"database/sql"
	"user/handlers"

	//"user/models
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
);	

var db *sql.DB
var dbConnected bool // Track DB connection status


func main() {
	var err error

	err = godotenv.Load("../../../.env")
    if err != nil {
        log.Fatalf("UserServiceErr: Error loading .env file: %v", err)
    }

	DB_AUTH := os.Getenv("DB_AUTH") 

	db, err = sql.Open("mysql", DB_AUTH + "usears_db")
	if err != nil {
		log.Printf("Initial database connection failed: %v", err)
	}

	dbConnected := make(chan bool)

	// Start a goroutine to retry the connection if it fails
	go retryDBConnection(DB_AUTH, dbConnected)
	
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", handlers.Home()) // This handler is available right away

	
	fmt.Println("Listening at port 5000")
	corsHandler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe("localhost:5000", corsHandler))
}

func retryDBConnection(dbAuth string, dbConnected chan bool) {
	for {
		if err := db.Ping(); err != nil {
			log.Printf("Database connection failed: %v. Retrying in 30 seconds...", err)

			// Try to reinitialize the database connection
			var err error
			db, err = sql.Open("mysql", dbAuth+"users_db")
			if err != nil {
				log.Printf("Error reinitializing database connection: %v", err)
			}
		} else {
			log.Println("Database connection successful!")
			dbConnected <- true
			close(dbConnected)
			SetDBConnected(true) // Update global dbConnected flag
			return
		}

		// Wait 1 minute before retrying
		time.Sleep(5 * time.Second)
	}
}

func SetDBConnected(status bool) {
	dbConnected = status
}

// HOME PAGE DOESNT DYNAMICALLY CHANGE WHEN THE DB CONNECTS SUCCESFFULY
