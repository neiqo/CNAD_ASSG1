package handlers

import (
	"fmt"
	"net/http"
)

func Home(getDBStatus func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !getDBStatus() {
			http.Error(w, "Error: User Service failed to connect to the database", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User Service connected to the database successfully!")
	}
}
