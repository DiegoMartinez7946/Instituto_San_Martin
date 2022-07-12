package services

import (
	"errors"
	"strings"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	jwt "github.com/dgrijalva/jwt-go"
)

/* G_UserEmail user will can access everywhere */
var GUserEmail string

/* G_UserID will can access everywhere */
var GUserID string

/* G_UserType will can access everywhere */
var GUserType string

/* ProcessToken verify our incoming token with the secret */
func ProcessToken(tk string) (*models.Claim, bool, string, error) {
	mySecret := []byte("diego_martinez")
	claims := &models.Claim {}

	splitToken := strings.Split(tk, "Bearer")
	if len(splitToken) != 2 {
		return claims, false, string(""), errors.New("formato de token invalido") // va sin signos este tipo de errores
	}

	tk = strings.TrimSpace(splitToken[1])

	// recibe el token, lo guarda en claims y el tercer parametro verifica el token con mySecret
	tkn, err := jwt.ParseWithClaims(tk, claims, func(token *jwt.Token)(interface{}, error) {
		return mySecret, nil
	})

	if err != nil {
		return claims, false, string(""), err
	}
	if !tkn.Valid {
		return claims, false, string(""), errors.New("token invalido")
	}

	_, find, _ := db.CheckExistUser(claims.Email)
	if find == true {
		GUserEmail = claims.Email
		GUserID = claims.ID.Hex()
		GUserType = claims.Type
	}

	return claims, find, GUserID, nil
}
