package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Teacher (docente) — colección teacher */
type Teacher struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name                string               `bson:"name" json:"name"`
	Email               string               `bson:"email" json:"email"`
	Phone               string               `bson:"phone" json:"phone"`
	DNI                 string               `bson:"dni" json:"dni"`
	Address             string               `bson:"address" json:"address"`
	EnseniaEn           []string             `bson:"enseniaen" json:"enseniaEn"`
	DegreeIDs           []primitive.ObjectID `bson:"degreeids" json:"degreeIds"`
	TituloHabilitanteID primitive.ObjectID   `bson:"titulohabilitanteid" json:"tituloHabilitanteId"`
	ModalidadID         primitive.ObjectID   `bson:"modalidadid" json:"modalidadId"`
	CreatedAt           time.Time            `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt           time.Time            `bson:"updatedat" json:"updatedAt,omitempty"`
}
