// Package models provides ...
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/***************************************************************/
/***************************************************************/
/* Shift model for the mongo DB */
type Shift struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:""`
	Type string `bson:"type" json:"type,omitempty"`
}
