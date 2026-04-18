// Package models provides ...
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/***************************************************************/
/***************************************************************/
/* Degree model for the mongo DB */
type Degree struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name,omitempty"`
	Nivel       string             `bson:"nivel" json:"nivel,omitempty"`
	StudyPlanID string             `bson:"studyplanid,omitempty" json:"studyPlanId,omitempty"` // numeroresolucion del plan (misma cadena que en studyplan)
	Turnos      []string           `bson:"turnos" json:"turnos"`
	Active      bool               `bson:"active" json:"active"`
}
