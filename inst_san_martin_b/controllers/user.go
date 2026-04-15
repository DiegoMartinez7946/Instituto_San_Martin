// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/***************************************************************/
/***************************************************************/
/* InsertUser insert one user role */
func InsertUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
		m := models.Response {
			Message: "Error en los parametros de entrada",
			Code: 400,
			Ok: false,
		}
		json.NewEncoder(w).Encode(m)
		return
	}
	
	if len(user.Email) == 0 {
		m := models.Response {
			Message: "El email de usuario es requerido",
			Code: 400,
			Ok: false,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	if len(user.Password) < 6 {
		m := models.Response {
			Message: "El password debe tener al menos 6 caracteres",
			Code: 400,
			Ok: false,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	msg, code, err := services.InsertUserService(user)
	if code != 201 {
		m := models.Response {
			Message: "Error al insertar el usuario " + msg,
			Code: code,
			Ok: false,
		}
		json.NewEncoder(w).Encode(m)
	  return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	m := models.Response {
		Message: "El usuario se registro correctamente",
		Code: code,
		Data: msg,
	}
	json.NewEncoder(w).Encode(m)
}

/***************************************************************/
/***************************************************************/
/* Login is the door to get in in the app */
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		m := models.Response {
			Message: "Faltan los parametros de entrada",
			Code: 400,
			Ok: false,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	u.Email = strings.TrimSpace(u.Email)
	if msg, ok := services.ValidateCorreoElectronicoRequired(u.Email); !ok {
		m := models.Response{
			Message: msg,
			Code:    400,
			Ok:      false,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	res, exists, err := services.LoginService(u)
	if exists == false || err!= nil {
		m := models.Response {
			Message: "Usuario y/o contrasena invalidos",
			Code: 400,
			Ok: false,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	expirationTime := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie {
		Name: "token",
		Value: res.Token,
		Expires: expirationTime,
	})
}

/***************************************************************/
/***************************************************************/
/* GetUsers get all the users */
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// call the services
	result, status, err := services.GetUsersService()
	if status == false {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		m := models.Response{
			Message: "No se puedo obtener la lista de usuarios" + errMsg,
			Code:    400,
			Ok:      false,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	m := models.Response {
		Message: "La lista de usuarios se ha obtenido exitosamente",
		Code: 200,
		Data: result,
	}
	json.NewEncoder(w).Encode(m)
}

/***************************************************************/
/***************************************************************/
/* ChangePassword allow us change the user password */
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var cp models.OldNewPassword

	err := json.NewDecoder(r.Body).Decode(&cp)
	if err != nil {
		m := models.Response {
			Message: "Parametros de entrada incorrectos " + err.Error(),
			Code: 400,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	result, status, err := services.ChangePasswordService(cp)
	if status == false {
		m := models.Response {
			Message: "Error al actualizar la password " + err.Error(),
			Code: 400,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

/***************************************************************/
/***************************************************************/
/* BlankPassword allow the admin change the password */
func BlankPassword(w http.ResponseWriter, r *http.Request) {
	var cp models.OldNewPassword

	err := json.NewDecoder(r.Body).Decode(&cp)
	if err != nil {
		m := models.Response {
			Message: "Parametros de entrada incorrectos",
			Code: 400,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	msg, code, err := services.BlankPasswordServices(cp)
	if err != nil || code != 200 {
		m := models.Response {
			Message: "Error al blanquear la password. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := models.Response {
		Message: msg,
		Code: code,
	}
	json.NewEncoder(w).Encode(res)
}

func mergeUserShiftSliceFromBody(ids []string, legacy string) []string {
	out := make([]string, 0, len(ids)+1)
	for _, x := range ids {
		if t := strings.TrimSpace(x); t != "" {
			out = append(out, t)
		}
	}
	if t := strings.TrimSpace(legacy); t != "" {
		out = append(out, t)
	}
	return out
}

type userUpdateBody struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	DNI       string   `json:"dni"`
	Address   string   `json:"address"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	UserType  string   `json:"userType"`
	ShiftID   string   `json:"shiftId"`   // legacy: un solo id
	ShiftIDs  []string `json:"shiftIds"` // preferido
	DegreeIDs []string `json:"degreeIds"`
	Modalidad string   `json:"modalidad"`
	Condicion string   `json:"condicion"`
	// Hex ObjectId de la ficha student cuando el listado la vinculó (más fiable que solo email).
	StudentRecordID string `json:"studentRecordId"`
}

/* UpdateUser PUT /user */
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var body userUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400, Ok: false})
		return
	}
	shiftIDs := mergeUserShiftSliceFromBody(body.ShiftIDs, body.ShiftID)
	msg, code, err := services.UpdateUserService(body.ID, body.Name, body.DNI, body.Address, body.Phone, body.Email, body.Password, body.UserType, shiftIDs, body.DegreeIDs, body.Modalidad, body.Condicion, body.StudentRecordID)
	if err != nil || code != 200 {
		_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code, Ok: false})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code, Ok: true})
}

type userActivePayload struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

/* ChangeActiveUser PUT /user/active */
func ChangeActiveUser(w http.ResponseWriter, r *http.Request) {
	var p userActivePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400})
		return
	}
	oid, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Id de usuario invalido", Code: 400})
		return
	}
	msg, code, err := services.UpdateUserActiveService(oid, p.Active)
	if err != nil || code != 200 {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Error al actualizar el estado. " + msg, Code: code})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
}
