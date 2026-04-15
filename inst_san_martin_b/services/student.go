package services

import (
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
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

/* InsertStudentService */
func InsertStudentService(s models.Student) (string, int, error) {
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

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	id, err := db.InsertStudentDB(s)
	if err != nil {
		return id, 400, err
	}
	return id, 201, nil
}

/* UpdateStudentService */
func UpdateStudentService(s models.Student) (string, int, error) {
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

	s.UpdatedAt = time.Now()
	okDB, err := db.UpdateStudentDB(s)
	if err != nil || !okDB {
		return "No se pudo actualizar el alumno", 400, err
	}
	return "El alumno se actualizo correctamente", 200, nil
}
