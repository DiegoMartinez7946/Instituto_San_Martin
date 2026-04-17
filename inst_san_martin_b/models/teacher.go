package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TeacherCareerAssignment título habilitante y modalidad por carrera (jerarquía implícita en el nivel de la carrera).
type TeacherCareerAssignment struct {
	DegreeID            primitive.ObjectID `bson:"degreeid" json:"degreeId"`
	TituloHabilitanteID primitive.ObjectID `bson:"titulohabilitanteid" json:"tituloHabilitanteId"`
	ModalidadID         primitive.ObjectID `bson:"modalidadid" json:"modalidadId"`
	ShiftID             primitive.ObjectID `bson:"shiftid" json:"shiftId"`
}

/* Teacher (docente) — colección teacher */
type Teacher struct {
	ID                primitive.ObjectID            `bson:"_id,omitempty" json:"id"`
	Name              string                          `bson:"name" json:"name"`
	Email             string                          `bson:"email" json:"email"`
	Phone             string                          `bson:"phone" json:"phone"`
	DNI               string                          `bson:"dni" json:"dni"`
	Address           string                          `bson:"address" json:"address"`
	EnseniaEn         []string                        `bson:"enseniaen" json:"enseniaEn"`
	Careers           []TeacherCareerAssignment       `bson:"careers" json:"careers"`
	LegacyDegreeIDs   []primitive.ObjectID            `bson:"degreeids,omitempty" json:"-"` // lectura de documentos viejos
	Active            bool                            `bson:"active" json:"active"`
	CreatedAt         time.Time                       `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt         time.Time                       `bson:"updatedat" json:"updatedAt,omitempty"`
}
