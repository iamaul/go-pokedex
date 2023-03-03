package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
}
