package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Email       string        `json:"email"`
	Password    string        `json:"password"`
	Name        string        `json:"name"`
	Surname     string        `json:"surname"`
	Birthdate   time.Time     `json:"birth_date"`
	Active      bool          `json:"active"`
	LastLoginAt time.Time     `json:"last_login_at"`
}
