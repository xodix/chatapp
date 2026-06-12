package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatMessage struct {
	ID      bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Author  bson.ObjectID `json:"author" bson:"author"`
	Message string        `json:"message" bson:"message"`
}

type Chat struct {
	ID       bson.ObjectID   `json:"id" bson:"_id,omitempty"`
	Messages []ChatMessage   `json:"messages" bson:"messages"`
	Members  []bson.ObjectID `json:"members" bson:"members"`
}
