// Package models provides ...
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* StudyPlan plan de estudio (colección studyplan). _id es el studyPlanId único. */
type StudyPlan struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre             string             `bson:"nombre" json:"nombre"`
	NumeroResolucion   string             `bson:"numeroresolucion" json:"numeroResolucion"`
	Cohorte            int                `bson:"cohorte" json:"cohorte"`                         // años de validez (0–9)
	Extension          int                `bson:"extension" json:"extension"`                     // años de prórroga (0–9)
	AnioAprobacion     int                `bson:"anioaprobacion" json:"anioAprobacion"`           // año calendario 4 dígitos
	AnioCaducidad      int                `bson:"aniocaducidad" json:"anioCaducidad"`             // anioAprobacion + cohorte + extension
	CargaHoraria       int                `bson:"cargahoraria" json:"cargaHoraria"`               // horas cátedra (1–99999)
	Active             bool               `bson:"active" json:"active"`
	CreatedAt          time.Time          `bson:"createdat,omitempty" json:"createdAt,omitempty"`
	UpdatedAt          time.Time          `bson:"updatedat,omitempty" json:"updatedAt,omitempty"`
}
