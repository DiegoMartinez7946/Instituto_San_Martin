package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type teacherPayload struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Email               string   `json:"email"`
	Phone               string   `json:"phone"`
	DNI                 string   `json:"dni"`
	Address             string   `json:"address"`
	EnseniaEn           []string `json:"enseniaEn"`
	DegreeIDs           []string `json:"degreeIds"`
	TituloHabilitanteID string   `json:"tituloHabilitanteId"`
	ModalidadID         string   `json:"modalidadId"`
}

func parseTeacherPayload(p teacherPayload) (models.Teacher, error) {
	var t models.Teacher
	t.Name = p.Name
	t.Email = p.Email
	t.Phone = p.Phone
	t.DNI = p.DNI
	t.Address = p.Address
	t.EnseniaEn = append([]string(nil), p.EnseniaEn...)

	if p.ID != "" {
		oid, err := primitive.ObjectIDFromHex(p.ID)
		if err != nil {
			return t, err
		}
		t.ID = oid
	}
	for _, hex := range p.DegreeIDs {
		oid, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			return t, err
		}
		t.DegreeIDs = append(t.DegreeIDs, oid)
	}
	if p.TituloHabilitanteID != "" {
		oid, err := primitive.ObjectIDFromHex(p.TituloHabilitanteID)
		if err != nil {
			return t, err
		}
		t.TituloHabilitanteID = oid
	}
	if p.ModalidadID != "" {
		oid, err := primitive.ObjectIDFromHex(p.ModalidadID)
		if err != nil {
			return t, err
		}
		t.ModalidadID = oid
	}
	return t, nil
}

/* GetTeachers */
func GetTeachers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, ok := services.GetTeachersService()
	if !ok {
		res := models.Response{Message: "Error al consultar los docentes", Code: 400}
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	res := models.Response{Data: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

/* InsertTeacher */
func InsertTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p teacherPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400})
		return
	}
	t, err := parseTeacherPayload(p)
	if err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Ids invalidos", Code: 400})
		return
	}
	msg, code, err := services.InsertTeacherService(t)
	if err != nil || code != 201 {
		res := models.Response{Message: "Error al insertar el docente. " + msg, Code: code}
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.Response{
		Message: "Se ha insertado el docente correctamente",
		Code:    code,
		Data:    msg,
	})
}

/* UpdateTeacher */
func UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p teacherPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400})
		return
	}
	t, err := parseTeacherPayload(p)
	if err != nil || t.ID.IsZero() {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Id de docente invalido", Code: 400})
		return
	}
	msg, code, err := services.UpdateTeacherService(t)
	if err != nil || code != 200 {
		res := models.Response{Message: "Error al actualizar el docente. " + msg, Code: code}
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
}
