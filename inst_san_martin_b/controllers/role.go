// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* GetRoles get all the user roles */
func GetRoles(w http.ResponseWriter, r *http.Request) {
	result, status := services.GetRolesService()
	if status == false {
		http.Error(w, "Error al consultar los roles", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

/* InsertRole insert one user role */
func InsertRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	err := json.NewDecoder(r.Body).Decode(&role)

	_, status, err := services.InsertRoleService(role)
	if err != nil {
		http.Error(w, "Ocurrio un error al registrar el rol " + err.Error(), 400)
		return
	}

	if status == false {
		http.Error(w, "No se ha logrado insertar el rol", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/* UpdateRole update one user role */
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	err := json.NewDecoder(r.Body).Decode(&role)

	var status bool
	status, err = services.UpdateRoleService(role)
	
	if err != nil {
		http.Error(w, "Ocurrio un error al intentar modificar el rol. Reintente nuevamente " + err.Error(), 400)
		return
	}

	if status == false {
		http.Error(w, "No se ha logrado modificar el rol", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/* DeleteRole delete one user role */
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		http.Error(w, "Debe enviar el id del rol para poder borrarlo", http.StatusBadRequest)
		return
	}

	err := services.DeleteRoleService(ID)
	if err != nil {
		http.Error(w, "Ocurrio un error al intentar borrar el rol" + err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
