package handlers

import (
	"fmt"
	"net/http"
)

func Status(getDBStatus func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !getDBStatus() {
			http.Error(w, "Error: Billing Service failed to connect to the database", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Billing Service connected to the database successfully!")
	}
}
