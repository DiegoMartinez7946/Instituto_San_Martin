package services

import (
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func shiftTypeToTurnoCode(typ string) string {
	u := strings.ToUpper(strings.TrimSpace(typ))
	switch u {
	case "MAÑANA", "MANANA":
		return "manana"
	case "TARDE":
		return "tarde"
	case "VESPERTINO":
		return "noche"
	default:
		return ""
	}
}

func degreeTurnosNorm(deg models.Degree) []string {
	var out []string
	for _, t := range deg.Turnos {
		c := strings.ToLower(strings.TrimSpace(t))
		if c != "" {
			out = append(out, c)
		}
	}
	return out
}

func ensureStudentDegreeShiftsDefaults(s *models.Student) {
	if s == nil || len(s.DegreeIDs) == 0 {
		s.DegreeShifts = nil
		return
	}
	if len(s.DegreeShifts) > 0 {
		return
	}
	sid, ok := db.GetFirstAvailableShiftID()
	if !ok {
		return
	}
	for _, did := range s.DegreeIDs {
		if did.IsZero() {
			continue
		}
		s.DegreeShifts = append(s.DegreeShifts, models.StudentDegreeShift{DegreeID: did, ShiftID: sid})
	}
}

func alignStudentDegreeShiftsWithDegreeOrder(s *models.Student) {
	if s == nil {
		return
	}
	byDeg := make(map[string]primitive.ObjectID)
	for _, row := range s.DegreeShifts {
		byDeg[row.DegreeID.Hex()] = row.ShiftID
	}
	var out []models.StudentDegreeShift
	for _, did := range s.DegreeIDs {
		if sid, ok := byDeg[did.Hex()]; ok {
			out = append(out, models.StudentDegreeShift{DegreeID: did, ShiftID: sid})
		}
	}
	s.DegreeShifts = out
}

func validateStudentDegreeShifts(s *models.Student) (string, bool) {
	if s == nil || len(s.DegreeIDs) == 0 {
		return "", true
	}
	if len(s.DegreeShifts) == 0 {
		return "Debe indicar el turno de cursada para cada carrera seleccionada", false
	}
	byDeg := make(map[primitive.ObjectID]primitive.ObjectID)
	for _, row := range s.DegreeShifts {
		if row.DegreeID.IsZero() || row.ShiftID.IsZero() {
			return "Cada carrera inscripta debe tener un turno valido", false
		}
		byDeg[row.DegreeID] = row.ShiftID
	}
	for _, did := range s.DegreeIDs {
		sid, ok := byDeg[did]
		if !ok {
			return "Debe indicar el turno de cursada para cada carrera seleccionada", false
		}
		ex, err := db.ShiftExistsByHexID(sid.Hex())
		if err != nil || !ex {
			return "Uno de los turnos indicados no es valido", false
		}
		sh, err := db.GetShiftDB(sid.Hex())
		if err != nil || sh.ID.IsZero() {
			return "Uno de los turnos indicados no es valido", false
		}
		deg, okD := db.GetDegreeByID(did)
		if !okD {
			return "Una o mas carreras no existen", false
		}
		allowed := degreeTurnosNorm(deg)
		if len(allowed) > 0 {
			code := shiftTypeToTurnoCode(sh.Type)
			okCode := false
			for _, a := range allowed {
				if a == code {
					okCode = true
					break
				}
			}
			if !okCode {
				return "El turno seleccionado no corresponde a los turnos habilitados para la carrera \"" + deg.Name + "\"", false
			}
		}
	}
	for did := range byDeg {
		found := false
		for _, x := range s.DegreeIDs {
			if x == did {
				found = true
				break
			}
		}
		if !found {
			return "Hay un turno asignado a una carrera que no esta seleccionada", false
		}
	}
	return "", true
}

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

	jmap := db.GetLevelMap()
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
	ensureStudentDegreeShiftsDefaults(&s)
	alignStudentDegreeShiftsWithDegreeOrder(&s)
	if msg, ok := validateStudentDegreeShifts(&s); !ok {
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
	ensureStudentDegreeShiftsDefaults(&s)
	alignStudentDegreeShiftsWithDegreeOrder(&s)
	if msg, ok := validateStudentDegreeShifts(&s); !ok {
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
