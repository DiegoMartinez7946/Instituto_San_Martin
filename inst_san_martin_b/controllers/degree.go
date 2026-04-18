// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* decodeDegreeJSON decodifica carrera; studyPlanId es el numeroresolucion del plan (cadena) o null. */
func decodeDegreeJSON(r *http.Request) (models.Degree, error) {
	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return models.Degree{}, err
	}
	var studyVal string
	var hadStudyPlanKey bool
	if v, ok := raw["studyPlanId"]; ok {
		hadStudyPlanKey = true
		if v == nil {
			studyVal = ""
		} else if s, ok := v.(string); ok {
			studyVal = strings.TrimSpace(s)
		}
	}
	delete(raw, "studyPlanId")
	delete(raw, "resolucionId")
	blob, err := json.Marshal(raw)
	if err != nil {
		return models.Degree{}, err
	}
	var d models.Degree
	if err := json.Unmarshal(blob, &d); err != nil {
		return models.Degree{}, err
	}
	if hadStudyPlanKey {
		d.StudyPlanID = studyVal
	}
	return d, nil
}

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
	degree, err := decodeDegreeJSON(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "JSON invalido: " + err.Error(),
			Code:    400,
		})
		return
	}

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
	degree, err := decodeDegreeJSON(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "JSON invalido: " + err.Error(),
			Code:    400,
		})
		return
	}

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
	if err := json.NewDecoder(r.Body).Decode(&degree); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "JSON invalido: " + err.Error(),
			Code:    400,
		})
		return
	}

	var code int
	var msg string
	var err error
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
