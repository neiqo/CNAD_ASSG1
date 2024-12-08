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

type Promotion struct {
	PromotionID  int     `json:"promotionID"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Discount     float64 `json:"discount"`
	IfPercentage bool    `json:"ifPercentage"`
}
