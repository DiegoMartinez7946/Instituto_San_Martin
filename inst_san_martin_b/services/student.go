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

func validateDegreeIDs(ids []primitive.ObjectID) (string, bool) {
	if len(ids) == 0 {
		return "Debe seleccionar al menos una carrera", false
	}
	for _, id := range ids {
		if id.IsZero() || !db.DegreeDocumentExists(id) {
			return "Una o mas carreras no existen", false
		}
	}
	return "", true
}

func validateEmail(email string) (string, bool) {
	e := strings.TrimSpace(email)
	if e == "" {
		return "", true
	}
	if !strings.Contains(e, "@") || len(e) < 5 {
		return "El correo electronico no es valido", false
	}
	return "", true
}

/* InsertStudentService */
func InsertStudentService(s models.Student) (string, int, error) {
	if len(s.Name) == 0 || len(s.DNI) == 0 {
		return "Nombre y DNI son obligatorios", 199, nil
	}
	if msg, ok := validateEmail(s.Email); !ok {
		return msg, 199, nil
	}
	_, exists := db.FindStudentByDNIDB(s.DNI, primitive.ObjectID{})
	if exists {
		return "Ya existe un alumno con ese DNI", 199, nil
	}
	msg, ok := validateDegreeIDs(s.DegreeIDs)
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
	if msg, ok := validateEmail(s.Email); !ok {
		return msg, 199, nil
	}
	_, exists := db.FindStudentByDNIDB(s.DNI, s.ID)
	if exists {
		return "Ya existe otro alumno con ese DNI", 199, nil
	}
	msg, ok := validateDegreeIDs(s.DegreeIDs)
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
