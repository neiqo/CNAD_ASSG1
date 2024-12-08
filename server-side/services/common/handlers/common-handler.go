package handlers

import (
	"common/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var db *sql.DB

func SetDBConnection(database *sql.DB) {
	db = database
}

func Status(getDBStatus func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !getDBStatus() {
			http.Error(w, "Error: Common Service failed to connect to the database", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Common Service connected to the database successfully!")
	}
}

func GetMemberBenefits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	membershipTier := vars["membershipTier"]

	membershipTier = strings.TrimSpace(membershipTier)

	rows, err := db.Query("SELECT benefitID, Name, Description FROM member_benefits WHERE membershipTier = ?", membershipTier)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error fetching member benefits. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var benefits []models.MemberBenefit
	for rows.Next() {
		var benefit models.MemberBenefit
		if err := rows.Scan(&benefit.BenefitID, &benefit.Name, &benefit.Description); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error scanning member benefit. Error: %v"}`, err), http.StatusInternalServerError)
			return
		}
		benefits = append(benefits, benefit)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(benefits) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No benefits found for the specified membership tier"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(benefits)
}

func GetPromotions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT promotionID, Name, Description, Discount, ifPercentage FROM promotions")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "Error fetching promotions. Error: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var promotions []models.Promotion
	for rows.Next() {
		var promo models.Promotion
		if err := rows.Scan(&promo.PromotionID, &promo.Name, &promo.Description, &promo.Discount, &promo.IfPercentage); err != nil {
			log.Printf("Error scanning promotion: %v", err)
			http.Error(w, fmt.Sprintf(`{"error": "Error scanning promotion. Error: %v"}`, err), http.StatusInternalServerError)
			return
		}
		promotions = append(promotions, promo)
	}

	log.Printf("Fetched %d promotions", len(promotions))

	w.Header().Set("Content-Type", "application/json")

	if len(promotions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No promotions found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(promotions)
}
