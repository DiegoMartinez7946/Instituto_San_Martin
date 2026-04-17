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

	GUserEmail, GUserID, GUserType = "", "", ""

	email := strings.TrimSpace(claims.Email)
	typ := strings.ToUpper(strings.TrimSpace(claims.Type))
	id := claims.ID

	if email != "" {
		user, inUser, _ := db.CheckExistUser(email)
		if inUser && user.Active {
			ur, ok, _ := db.GetUserDB(email)
			if ok && !ur.ID.IsZero() && ur.ID == id && strings.ToUpper(strings.TrimSpace(ur.Role)) == typ {
				GUserEmail = email
				GUserID = ur.ID.Hex()
				GUserType = typ
				return claims, true, GUserID, nil
			}
		}
	}

	if typ == "DOCENTE" && !id.IsZero() {
		t, ok := db.GetTeacherByID(id)
		if ok && t.Active && strings.EqualFold(strings.TrimSpace(t.Email), email) {
			GUserEmail = strings.TrimSpace(t.Email)
			GUserID = t.ID.Hex()
			GUserType = "DOCENTE"
			return claims, true, GUserID, nil
		}
	}

	if typ == "ALUMNO" && !id.IsZero() {
		st, ok := db.FindStudentByIDDB(id)
		if ok && st.Active && strings.EqualFold(strings.TrimSpace(st.Email), email) {
			GUserEmail = strings.TrimSpace(st.Email)
			GUserID = st.ID.Hex()
			GUserType = "ALUMNO"
			return claims, true, GUserID, nil
		}
	}

	return claims, false, "", nil
}
