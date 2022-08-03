// Package models provides ...
package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* User model for the mongo DB */
type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:""`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password,omitempty"`
	UserType string `bson:"usertype" json:"userType,omitempty"`
	CreatedAt time.Time `bson:"createdat" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedat" json:"updatedAt,omitempty"`
}

/* LoginResponse have the token which return in the login */
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

/* UserResponse get the user info */
type UserResponse struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:""`
	Email string `bson:"email" json:"email"`
	Role string `bson:"role" json:"role,omitempty"`
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

