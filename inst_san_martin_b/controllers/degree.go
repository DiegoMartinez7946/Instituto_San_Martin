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
/* GetDegrees get all the academy degrees */
func GetDegrees(w http.ResponseWriter, r *http.Request) {
	result, status := services.GetDegreesService()
	if status == false {
		res := models.Response {
			Message: "Error al consultar las carreras",
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
/* InsertDegree insert one academy degree */
func InsertDegree(w http.ResponseWriter, r *http.Request) {
	var degree models.Degree
	err := json.NewDecoder(r.Body).Decode(&degree)

	msg, code, err := services.InsertDegreeService(degree)
	if err != nil || code != 201 {
		res := models.Response {
			Message: "Error al insertar la carrera. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := models.Response {
		Message: "Se ha insertado la carrera correctamente",
		Code: code,
		Data: msg,
	}
	json.NewEncoder(w).Encode(res)
}

/***************************************************************/
/***************************************************************/
/* UpdateDegree update one academy degree */
func UpdateDegree(w http.ResponseWriter, r *http.Request) {
	var degree models.Degree
	err := json.NewDecoder(r.Body).Decode(&degree)

	var code int
	var msg string
	msg, code, err = services.UpdateDegreeService(degree)
	
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al actualizar la carrera. " + msg,
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
/* ChangeActiveDegree update status degree */
func ChangeActiveDegree(w http.ResponseWriter, r *http.Request) {
	var degree models.Degree
	err := json.NewDecoder(r.Body).Decode(&degree)

	var code int
	var msg string
	msg, code, err = services.UpdateDegreeStatusService(degree)
	
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al actualizar el estado de la carrera. " + msg,
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
