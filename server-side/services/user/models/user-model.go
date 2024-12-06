package models

type User struct {
	userID         int
	Name           string
	Email          string
	contactNo      int
	hashedPassword string
	membershipTier string
}
