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
	Modalidad      string               `bson:"modalidad,omitempty" json:"modalidad"`   // PRESENCIAL | REMOTO (json sin omitempty para que el front siempre reciba la clave)
	Condicion      string               `bson:"condicion,omitempty" json:"condicion"` // REGULAR | LIBRE
	DegreeIDs       []primitive.ObjectID `bson:"degreeids,omitempty" json:"degreeIds"`
	// Algunos documentos usan camelCase en MongoDB; se fusiona en lectura (db/student.go).
	DegreeIDsLegacy []primitive.ObjectID `bson:"degreeIds,omitempty" json:"-"`
	Active          bool                 `bson:"active" json:"active"`
	CreatedAt      time.Time            `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt      time.Time            `bson:"updatedat" json:"updatedAt,omitempty"`
}
