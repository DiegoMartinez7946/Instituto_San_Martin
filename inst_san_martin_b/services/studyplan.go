package services

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

var reNumeroResolucion = regexp.MustCompile(`^\d{4}/\d{2}$`)

func normalizeNumeroResolucion(s string) string {
	return strings.TrimSpace(s)
}

func validateNombrePlan(nombre string) (string, int, error) {
	n := strings.TrimSpace(nombre)
	if n == "" {
		return "El nombre del plan de estudio no puede estar vacio", 199, nil
	}
	return n, 0, nil
}

func validateStudyPlanFields(nombre, nr string, cohorte, extension, anioApr, carga int) (string, int, error) {
	if msg, code, err := validateNombrePlan(nombre); code == 199 {
		return msg, code, err
	}
	nr = normalizeNumeroResolucion(nr)
	if !reNumeroResolucion.MatchString(nr) {
		return "Numero de resolucion invalido: use el formato 6555/24 (cuatro digitos, barra, dos digitos de año)", 199, nil
	}
	if cohorte < 0 || cohorte > 9 {
		return "Cohorte (años de validez) debe ser un digito entre 0 y 9", 199, nil
	}
	if extension < 0 || extension > 9 {
		return "Extension debe ser un digito entre 0 y 9", 199, nil
	}
	if anioApr < 1000 || anioApr > 9999 {
		return "Año de aprobacion debe tener 4 digitos", 199, nil
	}
	suffix, err := strconv.Atoi(nr[len(nr)-2:])
	if err == nil && suffix != anioApr%100 {
		return "Los dos digitos tras la barra deben coincidir con el año de aprobacion (ej. 2024 -> .../24)", 199, nil
	}
	if carga < 1 || carga > 99999 {
		return "Carga horaria debe ser un entero entre 1 y 99999", 199, nil
	}
	return "", 0, nil
}

func computeAnioCaducidad(anioApr, cohorte, extension int) int {
	return anioApr + cohorte + extension
}

/* GetStudyPlansService lista planes de estudio */
func GetStudyPlansService() ([]*models.StudyPlan, bool) {
	return db.GetStudyPlansDB()
}

/* InsertStudyPlanService alta de plan */
func InsertStudyPlanService(s models.StudyPlan) (string, int, error) {
	if s.Active != true {
		return "El plan de estudio debe estar activo al crearse", 199, nil
	}
	msg, code, err := validateStudyPlanFields(s.Nombre, s.NumeroResolucion, s.Cohorte, s.Extension, s.AnioAprobacion, s.CargaHoraria)
	if code == 199 {
		return msg, code, err
	}
	nombreOK := strings.TrimSpace(s.Nombre)
	nr := normalizeNumeroResolucion(s.NumeroResolucion)
	if db.IsNumeroResolucionTakenByOther(nr, s.ID) {
		return "Ya existe un plan con ese numero de resolucion", 199, nil
	}
	anioCad := computeAnioCaducidad(s.AnioAprobacion, s.Cohorte, s.Extension)
	now := time.Now()
	row := models.StudyPlan{
		Nombre:           nombreOK,
		NumeroResolucion: nr,
		Cohorte:          s.Cohorte,
		Extension:        s.Extension,
		AnioAprobacion:   s.AnioAprobacion,
		AnioCaducidad:    anioCad,
		CargaHoraria:     s.CargaHoraria,
		Active:           true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	msgDB, err := db.InsertStudyPlanDB(row)
	if err != nil {
		return msgDB, 400, err
	}
	return msgDB, 201, nil
}

/* UpdateStudyPlanService actualiza plan */
func UpdateStudyPlanService(s models.StudyPlan) (string, int, error) {
	if s.ID.IsZero() {
		return "Id de plan de estudio invalido o no enviado", 199, nil
	}
	msg, code, err := validateStudyPlanFields(s.Nombre, s.NumeroResolucion, s.Cohorte, s.Extension, s.AnioAprobacion, s.CargaHoraria)
	if code == 199 {
		return msg, code, err
	}
	nombreOK := strings.TrimSpace(s.Nombre)
	nr := normalizeNumeroResolucion(s.NumeroResolucion)
	if db.IsNumeroResolucionTakenByOther(nr, s.ID) {
		return "Ya existe otro plan con ese numero de resolucion", 199, nil
	}
	s.Nombre = nombreOK
	s.NumeroResolucion = nr
	s.AnioCaducidad = computeAnioCaducidad(s.AnioAprobacion, s.Cohorte, s.Extension)
	s.UpdatedAt = time.Now()
	_, err = db.UpdateStudyPlanDB(s)
	if err != nil {
		if errors.Is(err, db.ErrStudyPlanNotFound) {
			return "No se encontro el plan de estudio con el id indicado", 199, nil
		}
		return "Hubo un error al actualizar el plan de estudio en la base", 400, err
	}
	return "El plan de estudio se actualizo correctamente", 200, nil
}

/* UpdateStudyPlanActiveService activo / inactivo */
func UpdateStudyPlanActiveService(s models.StudyPlan) (string, int, error) {
	if s.ID.IsZero() {
		return "Id de plan de estudio invalido o no enviado", 199, nil
	}
	_, err := db.UpdateStudyPlanActiveDB(s.ID, s.Active)
	if err != nil {
		if errors.Is(err, db.ErrStudyPlanNotFound) {
			return "No se encontro el plan de estudio con el id indicado", 199, nil
		}
		return "Hubo un error al actualizar el estado del plan de estudio en la base", 400, err
	}
	return "El estado del plan de estudio se actualizo correctamente", 200, nil
}
