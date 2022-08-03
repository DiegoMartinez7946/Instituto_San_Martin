// Package models provides ...
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/***************************************************************/
/***************************************************************/
/* TestType model for the mongo DB */
type TestType struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:""`
	Type string `bson:"type" json:"type,omitempty"`
}
