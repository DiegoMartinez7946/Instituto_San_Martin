package services

import (
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* GetStudentsService */
func GetStudentsService() ([]*models.Student, bool) {
	return db.GetStudentsDB()
}

func validateStudentDegreesEnrollment(s models.Student) (string, bool) {
	if len(s.DegreeIDs) == 0 {
		return "Debe seleccionar al menos una carrera", false
	}

	nivelAlumno := strings.ToLower(strings.TrimSpace(s.NivelAprobado))
	if nivelAlumno == "" {
		return "Debe indicar el nivel aprobado del alumno", false
	}

	jmap := db.GetJerarquiaMap()
	ha, ok := jmap[nivelAlumno]
	if !ok {
		return "Nivel aprobado no valido", false
	}

	for _, id := range s.DegreeIDs {
		if id.IsZero() || !db.DegreeDocumentExists(id) {
			return "Una o mas carreras no existen", false
		}
		deg, okDeg := db.GetDegreeByID(id)
		if !okDeg {
			return "Una o mas carreras no existen", false
		}
		niv := strings.ToLower(strings.TrimSpace(deg.Nivel))
		hn, okN := jmap[niv]
		if !okN || niv == "" {
			return "La carrera \"" + deg.Name + "\" no tiene un nivel valido configurado", false
		}
		if ha < hn-1 {
			return "El nivel aprobado del alumno no permite inscribirlo en la carrera \"" + deg.Name + "\"", false
		}
	}
	return "", true
}

func normalizeStudentModalidad(v string) string {
	u := strings.ToUpper(strings.TrimSpace(v))
	if u == "PRESENCIAL" || u == "REMOTO" {
		return u
	}
	return ""
}

func normalizeStudentCondicion(v string) string {
	u := strings.ToUpper(strings.TrimSpace(v))
	if u == "REGULAR" || u == "LIBRE" {
		return u
	}
	return ""
}

func validateStudentExtraModeFields(s *models.Student) (string, bool) {
	s.Modalidad = normalizeStudentModalidad(s.Modalidad)
	s.Condicion = normalizeStudentCondicion(s.Condicion)
	if !s.Active {
		// Regla pedida: al estar inactivo, los campos quedan en blanco.
		s.Modalidad = ""
		s.Condicion = ""
		return "", true
	}
	if s.Modalidad == "" {
		return "Debe indicar si el alumno es presencial o remoto", false
	}
	if s.Condicion == "" {
		return "Debe indicar si el alumno es regular o libre", false
	}
	return "", true
}

/* InsertStudentService newPortalPwd: contraseña de acceso al portal (opcional; mín. 6 caracteres si se informa). */
func InsertStudentService(s models.Student, newPortalPwd string) (string, int, error) {
	if len(s.Name) == 0 || len(s.DNI) == 0 {
		return "Nombre y DNI son obligatorios", 199, nil
	}
	if msg, ok := ValidateDNIAr(s.DNI); !ok {
		return msg, 199, nil
	}
	s.DNI = strings.TrimSpace(s.DNI)
	s.Email = strings.TrimSpace(s.Email)
	if msg, ok := ValidateCorreoElectronico(s.Email); !ok {
		return msg, 199, nil
	}
	s.Phone = NormalizeTelefonoDigits(s.Phone)
	if msg, ok := ValidateTelefono(s.Phone); !ok {
		return msg, 199, nil
	}
	_, exists := db.FindStudentByDNIDB(s.DNI, primitive.ObjectID{})
	if exists {
		return "Ya existe un alumno con ese DNI", 199, nil
	}
	s.NivelAprobado = strings.ToLower(strings.TrimSpace(s.NivelAprobado))
	msg, ok := validateStudentDegreesEnrollment(s)
	if !ok {
		return msg, 199, nil
	}
	if msg, ok := validateStudentExtraModeFields(&s); !ok {
		return msg, 199, nil
	}

	np := strings.TrimSpace(newPortalPwd)
	if np != "" {
		if len(np) < 6 {
			return "La contraseña de acceso al portal debe tener al menos 6 caracteres", 199, nil
		}
		h, errH := utils.EncryptPassword(np)
		if errH != nil {
			return "Error al cifrar la contraseña", 400, errH
		}
		s.Password = h
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	id, err := db.InsertStudentDB(s)
	if err != nil {
		return id, 400, err
	}
	return id, 201, nil
}

/* UpdateStudentService newPortalPwd: si no está vacío, actualiza contraseña de portal del alumno. */
func UpdateStudentService(s models.Student, newPortalPwd string) (string, int, error) {
	if s.ID.IsZero() {
		return "Falta el id del alumno", 199, nil
	}
	if len(s.Name) == 0 || len(s.DNI) == 0 {
		return "Nombre y DNI son obligatorios", 199, nil
	}
	if msg, ok := ValidateDNIAr(s.DNI); !ok {
		return msg, 199, nil
	}
	s.DNI = strings.TrimSpace(s.DNI)
	s.Email = strings.TrimSpace(s.Email)
	if msg, ok := ValidateCorreoElectronico(s.Email); !ok {
		return msg, 199, nil
	}
	s.Phone = NormalizeTelefonoDigits(s.Phone)
	if msg, ok := ValidateTelefono(s.Phone); !ok {
		return msg, 199, nil
	}
	_, exists := db.FindStudentByDNIDB(s.DNI, s.ID)
	if exists {
		return "Ya existe otro alumno con ese DNI", 199, nil
	}
	s.NivelAprobado = strings.ToLower(strings.TrimSpace(s.NivelAprobado))
	msg, ok := validateStudentDegreesEnrollment(s)
	if !ok {
		return msg, 199, nil
	}
	if msg, ok := validateStudentExtraModeFields(&s); !ok {
		return msg, 199, nil
	}

	s.UpdatedAt = time.Now()
	s.Password = ""
	okDB, err := db.UpdateStudentDB(s)
	if err != nil || !okDB {
		return "No se pudo actualizar el alumno", 400, err
	}
	if np := strings.TrimSpace(newPortalPwd); np != "" {
		if len(np) < 6 {
			return "La contraseña de acceso al portal debe tener al menos 6 caracteres", 199, nil
		}
		h, errH := utils.EncryptPassword(np)
		if errH != nil {
			return "Error al cifrar la contraseña", 400, errH
		}
		if errP := db.UpdateStudentPasswordHashDB(s.ID, h); errP != nil {
			return "Alumno actualizado pero no se pudo guardar la contraseña", 400, errP
		}
	}
	return "El alumno se actualizo correctamente", 200, nil
}

/* UpdateStudentActiveService solo activo/inactivo */
func UpdateStudentActiveService(id primitive.ObjectID, active bool) (string, int, error) {
	if id.IsZero() {
		return "Falta el id del alumno", 199, nil
	}
	okDB, err := db.UpdateStudentActiveDB(id, active)
	if err != nil || !okDB {
		return "No se pudo actualizar el estado del alumno", 400, err
	}
	if !active {
		_, _ = db.ClearStudentModeFieldsDB(id)
	}
	return "El estado del alumno se actualizo correctamente", 200, nil
}
