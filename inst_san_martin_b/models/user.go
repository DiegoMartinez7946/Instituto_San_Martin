// Package models provides ...
package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* User model for the mongo DB */
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password,omitempty"`
	UserType  string             `bson:"usertype" json:"userType,omitempty"`
	Name      string             `bson:"name" json:"name,omitempty"`
	DNI       string             `bson:"dni" json:"dni,omitempty"`
	Address   string             `bson:"address" json:"address,omitempty"`
	Phone     string             `bson:"phone" json:"phone,omitempty"`
	// shiftid legacy (un solo turno); shiftids es el campo actual
	ShiftID   string   `bson:"shiftid,omitempty" json:"shiftId,omitempty"`
	ShiftIDs []string `bson:"shiftids,omitempty" json:"shiftIds,omitempty"`
	// Carreras (ALUMNO / DOCENTE): ids Mongo como ObjectId; degreeIds en JSON es solo entrada API
	DegreeIDs     []primitive.ObjectID `bson:"degreeids,omitempty" json:"-"`
	DegreeIDHexes []string             `bson:"-" json:"degreeIds,omitempty"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedat" json:"updatedAt,omitempty"`
}

/* LoginResponse have the token which return in the login */
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

/* UserResponse get the user info */
type UserResponse struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:""`
	Email string             `bson:"email" json:"email"`
	Role  string             `bson:"role" json:"role,omitempty"`
}

/* UserListRow fila de listado de usuarios (administración) */
type UserListRow struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Name    string `json:"name"`
	DNI     string `json:"dni"`
	Address string `json:"address"`
	Phone     string `json:"phone"`
	ShiftIDs  []string `json:"shiftIds,omitempty"`
	ShiftType string   `json:"shiftType,omitempty"`
	DegreeIDs []string `json:"degreeIds,omitempty"`
	Active    bool   `json:"active"`
	// StudentLinked: cuenta user con rol ALUMNO y ficha en colección student (email o DNI).
	// Sin omitempty para que el front reciba siempre la clave (p. ej. studentLinked: false).
	StudentLinked bool   `json:"studentLinked"`
	// StudentRecordID hex de la ficha en colección student cuando hay vínculo (misma persona).
	StudentRecordID string `json:"studentRecordId"`
	Modalidad       string `json:"modalidad"`
	Condicion       string `json:"condicion"`
	// Source indica el origen del registro: vacío o "user" = colección user; "teacher" / "student" = ficha académica (no cuenta de login).
	Source string `json:"source,omitempty"`
}

/* Claim is the struct to process the jwt */
type Claim struct {
	Email string `json:"email"`
	ID primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Type string `json:"type"`
	jwt.StandardClaims
}

/* OldNewPassword is the struct to change the password */
type OldNewPassword struct {
	Email string `json:"email"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword string `json:"newPassword"`
}

