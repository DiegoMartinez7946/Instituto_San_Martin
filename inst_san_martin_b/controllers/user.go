// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
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
		m := models.Response {
			Message: "No se puedo obtener la lista de usuarios" + err.Error(),
			Code: 400,
			Ok: false,
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
