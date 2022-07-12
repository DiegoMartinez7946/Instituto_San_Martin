// Package models provides ...
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Role model for the mongo DB */
type Role struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:""`
	Type string `bson:"type" json:"type,omitempty"`
}
