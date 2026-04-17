package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type careerRowPayload struct {
	DegreeID            string `json:"degreeId"`
	TituloHabilitanteID string `json:"tituloHabilitanteId"`
	ModalidadID         string `json:"modalidadId"`
	ShiftID             string `json:"shiftId"`
}

type teacherPayload struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Phone       string             `json:"phone"`
	DNI         string             `json:"dni"`
	Address     string             `json:"address"`
	EnseniaEn   []string           `json:"enseniaEn"`
	Careers     []careerRowPayload `json:"careers"`
	Active      *bool              `json:"active"`
	NewPassword string             `json:"newPassword,omitempty"` // portal docente (administrativo)
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
	for _, row := range p.Careers {
		did, err := primitive.ObjectIDFromHex(row.DegreeID)
		if err != nil {
			return t, err
		}
		tid, err := primitive.ObjectIDFromHex(row.TituloHabilitanteID)
		if err != nil {
			return t, err
		}
		mid, err := primitive.ObjectIDFromHex(row.ModalidadID)
		if err != nil {
			return t, err
		}
		sid, err := primitive.ObjectIDFromHex(row.ShiftID)
		if err != nil {
			return t, err
		}
		t.Careers = append(t.Careers, models.TeacherCareerAssignment{
			DegreeID:            did,
			TituloHabilitanteID: tid,
			ModalidadID:         mid,
			ShiftID:             sid,
		})
	}
	if p.Active != nil {
		t.Active = *p.Active
	} else {
		t.Active = true
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
	msg, code, err := services.InsertTeacherService(t, strings.TrimSpace(p.NewPassword))
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
	msg, code, err := services.UpdateTeacherService(t, strings.TrimSpace(p.NewPassword))
	if err != nil || code != 200 {
		res := models.Response{Message: "Error al actualizar el docente. " + msg, Code: code}
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
}

type teacherActivePayload struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

/* ChangeActiveTeacher PUT /teacher/active */
func ChangeActiveTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p teacherActivePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400})
		return
	}
	oid, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Id de docente invalido", Code: 400})
		return
	}
	msg, code, err := services.UpdateTeacherActiveService(oid, p.Active)
	if err != nil || code != 200 {
		res := models.Response{Message: msg, Code: code}
		if code != 199 {
			res.Message = "Error al actualizar el estado del docente. " + msg
		}
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
}

type teacherPasswordPayload struct {
	ID          string `json:"id"`
	NewPassword string `json:"newPassword"`
}

/* ResetTeacherPassword PUT /teacher/password — solo administrativo (middleware). */
func ResetTeacherPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p teacherPasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		_ = json.NewEncoder(w).Encode(models.Response{Message: "Parametros invalidos", Code: 400})
		return
	}
	msg, code, err := services.ResetTeacherPasswordAdministrativo(p.ID, p.NewPassword)
	if err != nil || code != 200 {
		_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{Message: msg, Code: code})
}
