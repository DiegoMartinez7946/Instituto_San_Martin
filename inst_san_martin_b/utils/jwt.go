package utils

import (
	"time"

	"github.com/benjacifre10/san_martin_b/models"
	jwt "github.com/dgrijalva/jwt-go"
)

/* GenerateJWT generated the encrypt token with JWT */
func GenerateJWT(u models.UserResponse) (string, error) {
	mySecret := []byte("diego_martinez")

	payload := jwt.MapClaims {
		"email": u.Email,
		"_id": u.ID.Hex(),
		"type": u.Role.Type,
		"exp": time.Now().Add(time.Hour * 24).Unix(),// unix me lo transforma en un formato long
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString(mySecret)
	
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}
