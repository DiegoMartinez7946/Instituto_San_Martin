// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* GetStudyPlans GET /studyplan */
func GetStudyPlans(w http.ResponseWriter, r *http.Request) {
	result, status := services.GetStudyPlansService()
	if !status {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "Error al consultar los planes de estudio",
			Code:    400,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Data: result})
}

/* InsertStudyPlan POST /studyplan */
func InsertStudyPlan(w http.ResponseWriter, r *http.Request) {
	var studyPlan models.StudyPlan
	if err := json.NewDecoder(r.Body).Decode(&studyPlan); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "JSON invalido: " + err.Error(),
			Code:    400,
		})
		return
	}
	msg, code, err := services.InsertStudyPlanService(studyPlan)
	if err != nil || code != 201 {
		w.Header().Set("Content-Type", "application/json")
		status := http.StatusBadRequest
		if code == 199 {
			status = http.StatusUnprocessableEntity
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "Error al insertar el plan de estudio. " + msg,
			Code:    code,
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.Response{
		Message: "Se ha insertado el plan de estudio correctamente",
		Code:    code,
		Data:    msg,
	})
}

/* UpdateStudyPlan PUT /studyplan */
func UpdateStudyPlan(w http.ResponseWriter, r *http.Request) {
	var studyPlan models.StudyPlan
	if err := json.NewDecoder(r.Body).Decode(&studyPlan); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "JSON invalido: " + err.Error(),
			Code:    400,
		})
		return
	}
	msg, code, err := services.UpdateStudyPlanService(studyPlan)
	if err != nil || code != 200 {
		w.Header().Set("Content-Type", "application/json")
		status := http.StatusBadRequest
		if code == 199 {
			status = http.StatusUnprocessableEntity
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "Error al actualizar el plan de estudio. " + msg,
			Code:    code,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{
		Message: msg,
		Code:    code,
	})
}

/* ChangeActiveStudyPlan PUT /studyplan/active */
func ChangeActiveStudyPlan(w http.ResponseWriter, r *http.Request) {
	var studyPlan models.StudyPlan
	if err := json.NewDecoder(r.Body).Decode(&studyPlan); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "JSON invalido: " + err.Error(),
			Code:    400,
		})
		return
	}
	msg, code, err := services.UpdateStudyPlanActiveService(studyPlan)
	if err != nil || code != 200 {
		w.Header().Set("Content-Type", "application/json")
		status := http.StatusBadRequest
		if code == 199 {
			status = http.StatusUnprocessableEntity
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(models.Response{
			Message: "Error al actualizar el estado del plan de estudio. " + msg,
			Code:    code,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{
		Message: msg,
		Code:    code,
	})
}
