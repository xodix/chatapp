package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatMessage struct {
	Author  bson.ObjectID `json:"author" bson:"author"`
	Message string        `json:"message" bson:"message"`
	Name    string        `json:"name" bson:"name"`
	Surname string        `json:"surname" bson:"surname"`
}

type Chat struct {
	ID       bson.ObjectID   `json:"id" bson:"_id,omitempty"`
	Messages []ChatMessage   `json:"messages" bson:"messages"`
	Name     string          `json:"name" bson:"name"`
	Members  []bson.ObjectID `json:"members" bson:"members"`
}
