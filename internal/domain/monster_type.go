package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonsterType struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required,lte=4"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type MonsterTypeUpdate struct {
	ID   primitive.ObjectID `json:"_id,omitempty"`
	Name string             `json:"name"`
}

type MonsterTypeList struct {
	TotalCount   int            `json:"total_count"`
	TotalPages   int            `json:"total_pages"`
	Page         int            `json:"page"`
	Size         int            `json:"size"`
	HasMore      bool           `json:"has_more"`
	MonsterTypes []*MonsterType `json:"monster_types"`
}
