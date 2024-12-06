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
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var db *sql.DB
var dbConnected int32 // Track DB connection status

func main() {
	var err error

	err = godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("UserServiceErr: Error loading .env file: %v", err)
	}

	DB_AUTH := os.Getenv("DB_AUTH")

	db, err = sql.Open("mysql", DB_AUTH+"usears_db")
	if err != nil {
		log.Printf("Initial database connection failed: %v", err)
	}

	// Start a goroutine to retry the connection if it fails
	go retryDBConnection(DB_AUTH)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", handlers.Home(getDBStatus))

	fmt.Println("User Service listening at port 5010")
	corsHandler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe("localhost:5010", corsHandler))
}

func retryDBConnection(dbAuth string) {
	for {
		if err := db.Ping(); err != nil {
			log.Printf("Database connection failed: %v. Retrying in 5 seconds...", err)

			// Try to reinitialize the database connection
			var err error
			db, err = sql.Open("mysql", dbAuth+"users_db")
			if err != nil {
				log.Printf("Error reinitializing database connection: %v", err)
			}
		} else {
			log.Println("Database connection successful!")
			SetDBConnected(true)
			return
		}

		// Wait 1 minute before retrying
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

// HOME PAGE DOESNT DYNAMICALLY CHANGE WHEN THE DB CONNECTS SUCCESFFULY
