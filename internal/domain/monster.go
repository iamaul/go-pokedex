package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Monster struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	MonsterTypes []primitive.ObjectID `json:"monster_types" bson:"monster_types"`
	Name         string               `json:"name" bson:"name"`
	ImageUrl     string               `json:"image_url" bson:"image_url"`
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

type MonsterUpdate struct {
	ID          primitive.ObjectID `json:"_id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Size        float32            `json:"size"`
	Weight      float32            `json:"weight"`
	Hp          int32              `json:"hp"`
	Attack      int32              `json:"attack"`
	Defense     int32              `json:"defense"`
	Speed       int32              `json:"speed"`
}

type MonsterTypeBody struct {
	MonsterTypeID primitive.ObjectID `json:"monster_type_id"`
}

type MonsterList struct {
	TotalCount int        `json:"total_count"`
	TotalPages int        `json:"total_pages"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
	HasMore    bool       `json:"has_more"`
	Monsters   []*Monster `json:"monsters"`
}
