package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Monster struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	MonsterTypes []primitive.ObjectID `json:"monster_types" bson:"monster_types"`
	Name         string               `json:"name" bson:"name"`
	Description  string               `json:"description" bson:"description"`
	Size         float32              `json:"size" bson:"size"`
	Weight       float32              `json:"weight" bson:"weight"`
	Hp           int32                `json:"hp" bson:"hp"`
	Attack       int32                `json:"attack" bson:"attack"`
	Defense      int32                `json:"defense" bson:"defense"`
	Speed        int32                `json:"speed" bson:"speed"`
	CreatedAt    time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at" bson:"updated_at"`
}
