// Package controllers provides ...
package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* studentPayload decodes JSON degree ids as hex strings */
type degreeShiftPayload struct {
	DegreeID string `json:"degreeId"`
	ShiftID  string `json:"shiftId"`
}

type studentPayload struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Email         string               `json:"email"`
	Phone         string               `json:"phone"`
	DNI           string               `json:"dni"`
	Address       string               `json:"address"`
	NivelAprobado string               `json:"nivelAprobado"`
	Modalidad     string               `json:"modalidad"`
	Condicion     string               `json:"condicion"`
	DegreeIDs     []string             `json:"degreeIds"`
	DegreeShifts  []degreeShiftPayload `json:"degreeShifts"`
	Active        *bool                `json:"active"`
	NewPassword   string               `json:"newPassword,omitempty"` // portal alumno (administrativo)
}

func parseStudentPayload(p studentPayload) (models.Student, error) {
	var s models.Student
	s.Name = p.Name
	s.Email = p.Email
	s.Phone = p.Phone
	s.DNI = p.DNI
	s.Address = p.Address
	s.NivelAprobado = p.NivelAprobado
	s.Modalidad = p.Modalidad
	s.Condicion = p.Condicion

	if p.ID != "" {
		oid, err := primitive.ObjectIDFromHex(p.ID)
		if err != nil {
			return s, err
		}
		s.ID = oid
	}

	for _, hex := range p.DegreeIDs {
		oid, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			return s, err
		}
		s.DegreeIDs = append(s.DegreeIDs, oid)
	}
	for _, row := range p.DegreeShifts {
		did, err := primitive.ObjectIDFromHex(strings.TrimSpace(row.DegreeID))
		if err != nil {
			return s, err
		}
		sid, err := primitive.ObjectIDFromHex(strings.TrimSpace(row.ShiftID))
		if err != nil {
			return s, err
		}
		s.DegreeShifts = append(s.DegreeShifts, models.StudentDegreeShift{DegreeID: did, ShiftID: sid})
	}
	if p.Active != nil {
		s.Active = *p.Active
	} else {
		s.Active = true
	}
	return s, nil
}

/* GetStudents */
func GetStudents(w http.ResponseWriter, r *http.Request) {
	result, status := services.GetStudentsService()
	if !status {
		res := models.Response{
			Message: "Error al consultar los alumnos",
			Code:    400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := models.Response{
		Data: result,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

/* InsertStudent */
func InsertStudent(w http.ResponseWriter, r *http.Request) {
	var p studentPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		res := models.Response{
			Message: "Parametros invalidos",
			Code:    400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	s, err := parseStudentPayload(p)
	if err != nil {
		res := models.Response{
			Message: "Ids invalidos",
			Code:    400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	msg, code, err := services.InsertStudentService(s, strings.TrimSpace(p.NewPassword))
	if err != nil || code != 201 {
		res := models.Response{
			Message: "Error al insertar el alumno. " + msg,
			Code:    code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := models.Response{
		Message: "Se ha insertado el alumno correctamente",
		Code:    code,
		Data:    msg,
	}
	json.NewEncoder(w).Encode(res)
}

/* UpdateStudent */
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var p studentPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		res := models.Response{
			Message: "Parametros invalidos",
			Code:    400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	s, err := parseStudentPayload(p)
	if err != nil || s.ID.IsZero() {
		res := models.Response{
			Message: "Id de alumno invalido",
			Code:    400,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	msg, code, err := services.UpdateStudentService(s, strings.TrimSpace(p.NewPassword))
	if err != nil || code != 200 {
		res := models.Response{
			Message: "Error al actualizar el alumno. " + msg,
			Code:    code,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := models.Response{
		Message: msg,
		Code:    code,
	}
	json.NewEncoder(w).Encode(res)
}

type studentPasswordPayload struct {
	ID          string `json:"id"`
	NewPassword string `json:"newPassword"`
}

/* ResetStudentPassword PUT /student/password — solo administrativo (middleware). */
func ResetStudentPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p studentPasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400})
		return
	}
	msg, code, err := services.ResetStudentPasswordAdministrativo(p.ID, p.NewPassword)
	if err != nil || code != 200 {
		_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
}

type studentActivePayload struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

/* ChangeActiveStudent PUT /student/active */
func ChangeActiveStudent(w http.ResponseWriter, r *http.Request) {
	var p studentActivePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400})
		return
	}
	oid, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Id de alumno invalido", Code: 400})
		return
	}
	msg, code, err := services.UpdateStudentActiveService(oid, p.Active)
	if err != nil || code != 200 {
		res := models.Response{Message: "Error al actualizar el estado del alumno. " + msg, Code: code}
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
}
