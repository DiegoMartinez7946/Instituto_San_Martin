// Package models provides ...
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Student model for the mongo DB */
type Student struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name           string               `bson:"name" json:"name"`
	Email          string               `bson:"email" json:"email"`
	Phone          string               `bson:"phone" json:"phone"`
	DNI            string               `bson:"dni" json:"dni"`
	Address        string               `bson:"address" json:"address"`
	NivelAprobado  string               `bson:"nivelaprobado" json:"nivelAprobado"`
	DegreeIDs      []primitive.ObjectID `bson:"degreeids" json:"degreeIds"`
	CreatedAt      time.Time            `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt      time.Time            `bson:"updatedat" json:"updatedAt,omitempty"`
}
