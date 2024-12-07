package models

import "time"

type Payment struct {
	PaymentID   int       `json:"paymentID"`
	UserID      int       `json:"userID"`
	BookingID   int       `json:"bookingID"`
	Status      string    `json:"status"`
	PromotionID int       `json:"promotionID"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"createdAt"`
}
