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
/* GetShifts get all the shifts of the academy */
func GetShifts(w http.ResponseWriter, r *http.Request) {
	result, code, err := services.GetShiftsService()
	if code != 200 {
		res := models.Response {
			Message: "Error al consultar los turnos " + err.Error(),
			Code: code,
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
/* InsertShift insert one academy shift */
func InsertShift(w http.ResponseWriter, r *http.Request) {
	var shift models.Shift
	err := json.NewDecoder(r.Body).Decode(&shift)

	msg, code, err := services.InsertShiftService(shift)
	if err != nil || code != 201 {
		res := models.Response {
			Message: "Error al insertar el turno. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := models.Response {
		Message: "Se ha insertado el turno correctamente",
		Code: code,
		Data: msg,
	}
	json.NewEncoder(w).Encode(res)
	w.Header().Set("Content-Type", "application/json")
}

/***************************************************************/
/***************************************************************/
/* UpdateShift update one academy shift */
func UpdateShift(w http.ResponseWriter, r *http.Request) {
	var shift models.Shift
	err := json.NewDecoder(r.Body).Decode(&shift)

	var code int
	var msg string
	msg, code, err = services.UpdateShiftService(shift)
	
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al actualizar el turno. " + msg,
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
/* DeleteShift delete one academy shift */
func DeleteShift(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		res := models.Response {
			Message: "Falta un parametro para borrar el turno",
			Code: 400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	msg, code, err := services.DeleteShiftService(ID)
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al borrar el turno. " + msg,
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
