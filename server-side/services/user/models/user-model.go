package models

type User struct {
	UserID         int    `json:"user_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ContactNo      string `json:"contact_no"`
	HashedPassword string `json:"hashed_password"`
	MembershipTier string `json:"membership_tier"`
}
