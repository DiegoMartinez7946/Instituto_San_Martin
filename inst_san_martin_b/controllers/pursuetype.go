// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/***************************************************************/
/***************************************************************/
/* GetPursueTypes get all the pursue types */
func GetPursueTypes(w http.ResponseWriter, r *http.Request) {
	result, status := services.GetPursueTypesService()
	if status == false {
		res := models.Response {
			Message: "Error al consultar las modalidades de cursado",
			Code: 400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := models.Response {
		Data: result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

/***************************************************************/
/***************************************************************/
/* InsertPursueType insert one pursue type */
func InsertPursueType(w http.ResponseWriter, r *http.Request) {
	var pursueType models.PursueType
	err := json.NewDecoder(r.Body).Decode(&pursueType)

	msg, code, err := services.InsertPursueTypeService(pursueType)
	if err != nil || code != 201 {
		res := models.Response {
			Message: "Error al insertar la modalidad de cursado. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := models.Response {
		Message: "Se ha insertado la modalidad de cursado correctamente",
		Code: code,
		Data: msg,
	}
	json.NewEncoder(w).Encode(res)
}

/***************************************************************/
/***************************************************************/
/* UpdatePursueType update one pursue type */
func UpdatePursueType(w http.ResponseWriter, r *http.Request) {
	var pursueType models.PursueType
	err := json.NewDecoder(r.Body).Decode(&pursueType)

	var code int
	var msg string
	msg, code, err = services.UpdatePursueTypeService(pursueType)
	
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al actualizar la modalidad de cursado. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := models.Response {
		Message: msg,
		Code: code,
	}
	json.NewEncoder(w).Encode(res)
}

/***************************************************************/
/***************************************************************/
/* DeletePursueType delete one pursue type */
func DeletePursueType(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		res := models.Response {
			Message: "Falta un parametro para borrar la modalidad de cursado",
			Code: 400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	msg, code, err := services.DeletePursueTypeService(ID)
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al borrar la modalidad de cursado. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := models.Response {
		Message: msg,
		Code: code,
	}
	json.NewEncoder(w).Encode(res)
}
