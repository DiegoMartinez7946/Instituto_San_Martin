// Package models provides ...
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* StudentDegreeShift turno de cursada por carrera (referencia a colección shift). */
type StudentDegreeShift struct {
	DegreeID primitive.ObjectID `bson:"degreeid" json:"degreeId"`
	ShiftID  primitive.ObjectID `bson:"shiftid" json:"shiftId"`
}

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
	// DegreeIDs se expone en JSON; en Mongo la lista de carreras sale de degreeshifts (degreeids en BSON solo lectura legacy).
	DegreeIDs       []primitive.ObjectID `bson:"degreeids,omitempty" json:"degreeIds"`
	DegreeShifts []StudentDegreeShift `bson:"degreeshifts,omitempty" json:"degreeShifts"`
	// Algunos documentos usan camelCase en MongoDB; se fusiona en lectura (db/student.go).
	DegreeIDsLegacy []primitive.ObjectID `bson:"degreeIds,omitempty" json:"-"`
	Active          bool                 `bson:"active" json:"active"`
	Password        string               `bson:"password,omitempty" json:"-"` // bcrypt; login portal alumno
	CreatedAt      time.Time            `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt      time.Time            `bson:"updatedat" json:"updatedAt,omitempty"`
}
