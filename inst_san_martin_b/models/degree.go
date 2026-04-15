// Package models provides ...
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/***************************************************************/
/***************************************************************/
/* Degree model for the mongo DB */
type Degree struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name,omitempty"`
	Nivel        string             `bson:"nivel" json:"nivel,omitempty"`
	ResolucionID string             `bson:"resolucionid" json:"resolucionId,omitempty"`
	Active       bool               `bson:"active" json:"active"`
}
