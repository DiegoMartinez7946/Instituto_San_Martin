package services

import (
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var allowedEnseniaNiveles = map[string]struct{}{
	"inicial":       {},
	"primario":      {},
	"secundario":    {},
	"terciario":     {},
	"universitario": {},
}

func normalizeEnseniaEn(list []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, raw := range list {
		k := strings.ToLower(strings.TrimSpace(raw))
		if k == "" {
			continue
		}
		if _, ok := allowedEnseniaNiveles[k]; !ok {
			continue
		}
		if _, dup := seen[k]; dup {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, k)
	}
	return out
}

func validateTeacherDegrees(ids []primitive.ObjectID) (string, bool) {
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

func validateTitularModalidad(t models.Teacher) (string, bool) {
	mod, ok := db.GetModalidadByID(t.ModalidadID)
	if !ok {
		return "Modalidad no valida", false
	}
	tit, ok := db.GetTituloHabilitanteByID(t.TituloHabilitanteID)
	if !ok {
		return "Titulo habilitante no valido", false
	}
	if strings.ToUpper(strings.TrimSpace(mod.Codigo)) == "TITULAR" &&
		strings.ToUpper(strings.TrimSpace(tit.Codigo)) != "SI" {
		return "Solo puede ser titular si tiene titulo habilitante SI", false
	}
	return "", true
}

/* GetTeachersService */
func GetTeachersService() ([]*models.Teacher, bool) {
	return db.GetTeachersDB()
}

/* InsertTeacherService */
func InsertTeacherService(t models.Teacher) (string, int, error) {
	if len(t.Name) == 0 || len(t.DNI) == 0 {
		return "Nombre y DNI son obligatorios", 199, nil
	}
	if msg, ok := ValidateDNIAr(t.DNI); !ok {
		return msg, 199, nil
	}
	t.DNI = strings.TrimSpace(t.DNI)
	t.Email = strings.TrimSpace(t.Email)
	if msg, ok := ValidateCorreoElectronico(t.Email); !ok {
		return msg, 199, nil
	}
	t.Phone = NormalizeTelefonoDigits(t.Phone)
	if msg, ok := ValidateTelefono(t.Phone); !ok {
		return msg, 199, nil
	}
	_, exists := db.FindTeacherByDNIDB(t.DNI, primitive.ObjectID{})
	if exists {
		return "Ya existe un docente con ese DNI", 199, nil
	}

	t.EnseniaEn = normalizeEnseniaEn(t.EnseniaEn)
	if len(t.EnseniaEn) == 0 {
		return "Debe indicar al menos un nivel en que ensenia", 199, nil
	}
	if msg, ok := validateTeacherDegrees(t.DegreeIDs); !ok {
		return msg, 199, nil
	}
	if t.TituloHabilitanteID.IsZero() || t.ModalidadID.IsZero() {
		return "Debe seleccionar titulo habilitante y modalidad", 199, nil
	}
	if msg, ok := validateTitularModalidad(t); !ok {
		return msg, 199, nil
	}

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	id, err := db.InsertTeacherDB(t)
	if err != nil {
		return id, 400, err
	}
	return id, 201, nil
}

/* UpdateTeacherService */
func UpdateTeacherService(t models.Teacher) (string, int, error) {
	if t.ID.IsZero() {
		return "Falta el id del docente", 199, nil
	}
	if len(t.Name) == 0 || len(t.DNI) == 0 {
		return "Nombre y DNI son obligatorios", 199, nil
	}
	if msg, ok := ValidateDNIAr(t.DNI); !ok {
		return msg, 199, nil
	}
	t.DNI = strings.TrimSpace(t.DNI)
	t.Email = strings.TrimSpace(t.Email)
	if msg, ok := ValidateCorreoElectronico(t.Email); !ok {
		return msg, 199, nil
	}
	t.Phone = NormalizeTelefonoDigits(t.Phone)
	if msg, ok := ValidateTelefono(t.Phone); !ok {
		return msg, 199, nil
	}
	_, exists := db.FindTeacherByDNIDB(t.DNI, t.ID)
	if exists {
		return "Ya existe otro docente con ese DNI", 199, nil
	}

	t.EnseniaEn = normalizeEnseniaEn(t.EnseniaEn)
	if len(t.EnseniaEn) == 0 {
		return "Debe indicar al menos un nivel en que ensenia", 199, nil
	}
	if msg, ok := validateTeacherDegrees(t.DegreeIDs); !ok {
		return msg, 199, nil
	}
	if t.TituloHabilitanteID.IsZero() || t.ModalidadID.IsZero() {
		return "Debe seleccionar titulo habilitante y modalidad", 199, nil
	}
	if msg, ok := validateTitularModalidad(t); !ok {
		return msg, 199, nil
	}

	t.UpdatedAt = time.Now()
	okDB, err := db.UpdateTeacherDB(t)
	if err != nil || !okDB {
		return "No se pudo actualizar el docente", 400, err
	}
	return "El docente se actualizo correctamente", 200, nil
}
