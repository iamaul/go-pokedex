package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonsterType struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Type      string             `json:"type" bson:"type"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
