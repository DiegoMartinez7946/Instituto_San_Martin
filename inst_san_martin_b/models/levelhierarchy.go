package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// LevelHierarchy maps an academic level to its hierarchy rank (Mongo collection "jerarquia").
type LevelHierarchy struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nivel     string             `bson:"nivel" json:"nivel"`
	Jerarquia int                `bson:"jerarquia" json:"jerarquia"`
}
