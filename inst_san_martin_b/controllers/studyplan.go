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
/* GetStudyPlans get all the study plans */
func GetStudyPlans(w http.ResponseWriter, r *http.Request) {
	result, status := services.GetStudyPlansService()
	if status == false {
		res := models.Response {
			Message: "Error al consultar los planes de estudio",
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
/* InsertStudyPlan insert one study plan */
func InsertStudyPlan(w http.ResponseWriter, r *http.Request) {
	var studyPlan models.StudyPlan
	err := json.NewDecoder(r.Body).Decode(&studyPlan)

	msg, code, err := services.InsertStudyPlanService(studyPlan)
	if err != nil || code != 201 {
		res := models.Response {
			Message: "Error al insertar el plan de estudio. " + msg,
			Code: code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := models.Response {
		Message: "Se ha insertado el plan de estudio correctamente",
		Code: code,
		Data: msg,
	}
	json.NewEncoder(w).Encode(res)
}

/***************************************************************/
/***************************************************************/
/* UpdateStudyPlan update one study plan */
func UpdateStudyPlan(w http.ResponseWriter, r *http.Request) {
	var studyPlan models.StudyPlan
	err := json.NewDecoder(r.Body).Decode(&studyPlan)

	var code int
	var msg string
	msg, code, err = services.UpdateStudyPlanService(studyPlan)
	
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al actualizar el plan de estudio. " + msg,
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
/* ChangeStateStudyPlan update status study plan */
func ChangeStateStudyPlan(w http.ResponseWriter, r *http.Request) {
	var studyPlan models.StudyPlan
	err := json.NewDecoder(r.Body).Decode(&studyPlan)

	var code int
	var msg string
	msg, code, err = services.UpdateStudyPlanStateService(studyPlan)
	
	if err != nil || code != 200 {
		res := models.Response {
			Message: "Error al actualizar el estado del plan de estudio. " + msg,
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
