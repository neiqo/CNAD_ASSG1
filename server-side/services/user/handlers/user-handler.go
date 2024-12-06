package handlers

import (
	"fmt"
	"net/http"
)

	var dbConnected = false // Assume initially disconnected

	func SetDBConnected(status bool) {
		dbConnected = status
	}

	func Home() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !dbConnected {
				http.Error(w, "Error: User Service failed to connect to the database", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "User Service connected to the database successfully!")
		}
	}