// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* InsertUser insert one user role */
func InsertUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
		http.Error(w, "Error en los datos recibidos " + err.Error(), 400)
		return
	}
	
	if len(user.Email) == 0 {
		http.Error(w, "El email de usuario es requerido", 400)
		return
	}

	if len(user.Password) < 6 {
		http.Error(w, "El password debe tener al menos 6 caracteres", 400)
		return
	}

	message, status, err := services.InsertUserService(user)
	if message != "" {
		http.Error(w, message, 400)
	  return
	}

	if err != nil {
		http.Error(w, "Ocurrio un error al insertar el usuario " + err.Error(), 400)
	  return
	}

	if status == false {
		http.Error(w, "No se ha logrado insertar el usuario", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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

/* GetUsers get all the users */
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// call the services
	result, status := services.GetUsersService()
	if status == false {
		http.Error(w, "Error al consultar los usuarios", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

/* ChangePassword allow us change the user password */
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var cp models.OldNewPassword

	err := json.NewDecoder(r.Body).Decode(&cp)
	if err != nil {
		http.Error(w, "Parametros de entrada incorrectos " + err.Error(), 400)
		return
	}

	result, status, err := services.ChangePasswordService(cp)
	if status == false {
		http.Error(w, "Error al actualizar la password", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
