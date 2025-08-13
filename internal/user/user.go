package user

import "time"

type User struct {
	ID           int       `json:"id"`
	PhoneNumber  string    `json:"phone_number"`
	RegisteredAt time.Time `json:"registered_at"`
}
