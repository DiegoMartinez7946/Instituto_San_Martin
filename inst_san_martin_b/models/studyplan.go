// Package models provides ...
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/***************************************************************/
/***************************************************************/
/* StudyPlan model for the mongo DB */
type StudyPlan struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:""`
	Name string `bson:"name" json:"name,omitempty"`
	Code string `bson:"code" json:"code,omitempty"`
	DegreeId string `bson:"degreeid" json:"degreeId,omitempty"`
	State bool `bson:"state" json:"state"`
	CreatedAt time.Time `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedat" json:"updatedAt,omitempty"`
}
