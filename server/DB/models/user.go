package models

import (
	"time"
)

type User struct {
	ID          *string   `json:"id" bson:"_id,omitempty"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Birthdate   time.Time `json:"birth_date"`
	Active      bool      `json:"active"`
	LastLoginAt time.Time `json:"last_login_at"`
}
