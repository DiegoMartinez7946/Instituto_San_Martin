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
/* GetTestTypes get all the test types */
func GetTestTypes(w http.ResponseWriter, r *http.Request) {
	result, status := services.GetTestTypesService()
	if status == false {
		res := models.Response {
			Message: "Error al consultar los tipos de examenes",
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
/* InsertTestType insert one test type */
func InsertTestType(w http.ResponseWriter, r *http.Request) {
	var testType models.TestType
	err := json.NewDecoder(r.Body).Decode(&testType)

	msg, code, err := services.InsertTestTypeService(testType)
	if err != nil || code != 201 {
		res := models.Response {
			Message: "Error al insertar el tipo de examen. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := models.Response {
		Message: "Se ha insertado el tipo de examen correctamente",
		Code: code,
		Data: msg,
	}
	json.NewEncoder(w).Encode(res)
}

/***************************************************************/
/***************************************************************/
/* UpdateTestType update one test type */
func UpdateTestType(w http.ResponseWriter, r *http.Request) {
	var testType models.TestType
	err := json.NewDecoder(r.Body).Decode(&testType)

	var code int
	var msg string
	msg, code, err = services.UpdateTestTypeService(testType)
	
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al actualizar el tipo de examen. " + msg,
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
/* DeleteTestType delete one test type */
func DeleteTestType(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		res := models.Response {
			Message: "Falta un parametro para borrar el tipo de examen",
			Code: 400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	msg, code, err := services.DeleteTestTypeService(ID)
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al borrar el tipo de examen. " + msg,
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
