package services

import (
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
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

func codigoTitulo(id primitive.ObjectID) (string, bool) {
	t, ok := db.GetTituloHabilitanteByID(id)
	if !ok {
		return "", false
	}
	return strings.ToUpper(strings.TrimSpace(t.Codigo)), true
}

func codigoModalidad(id primitive.ObjectID) (string, bool) {
	m, ok := db.GetModalidadByID(id)
	if !ok {
		return "", false
	}
	return strings.ToUpper(strings.TrimSpace(m.Codigo)), true
}

func validateTeacherCareers(t *models.Teacher) (string, bool) {
	if len(t.Careers) == 0 {
		return "Debe seleccionar al menos una carrera y definir titulo y modalidad por carrera", false
	}
	ens := make(map[string]struct{})
	for _, n := range t.EnseniaEn {
		ens[strings.ToLower(strings.TrimSpace(n))] = struct{}{}
	}
	seenDeg := make(map[string]struct{})
	for _, c := range t.Careers {
		if c.DegreeID.IsZero() {
			return "Carrera invalida en la lista", false
		}
		if !db.DegreeDocumentExists(c.DegreeID) {
			return "Una o mas carreras no existen", false
		}
		h := c.DegreeID.Hex()
		if _, dup := seenDeg[h]; dup {
			return "No puede repetir la misma carrera", false
		}
		seenDeg[h] = struct{}{}

		deg, ok := db.GetDegreeByID(c.DegreeID)
		if !ok {
			return "Una o mas carreras no existen", false
		}
		niv := strings.ToLower(strings.TrimSpace(deg.Nivel))
		if niv == "" {
			return "La carrera \"" + deg.Name + "\" no tiene nivel configurado", false
		}
		if _, ok := ens[niv]; !ok {
			return "La carrera \"" + deg.Name + "\" no esta cubierta por los niveles en que ensenia", false
		}

		if c.TituloHabilitanteID.IsZero() || c.ModalidadID.IsZero() {
			return "Cada carrera requiere titulo habilitante y modalidad", false
		}
		if c.ShiftID.IsZero() {
			return "Cada carrera requiere un turno", false
		}
		if _, ok := db.GetTituloHabilitanteByID(c.TituloHabilitanteID); !ok {
			return "Titulo habilitante no valido", false
		}
		if _, ok := db.GetModalidadByID(c.ModalidadID); !ok {
			return "Modalidad no valida", false
		}
		if exists, err := db.ShiftExistsByHexID(c.ShiftID.Hex()); err != nil || !exists {
			return "Turno no valido", false
		}

		tit, _ := codigoTitulo(c.TituloHabilitanteID)
		mod, _ := codigoModalidad(c.ModalidadID)
		if mod == "TITULAR" && tit != "SI" {
			return "En la carrera \"" + deg.Name + "\" solo puede ser titular con titulo habilitante SI", false
		}
		if tit != "SI" && mod != "PROVISIONAL" && mod != "SUPLENTE" {
			return "En la carrera \"" + deg.Name + "\" sin titulo habilitante solo puede modalidad provisional o suplente", false
		}
	}
	return "", true
}

/* GetTeachersService */
func GetTeachersService() ([]*models.Teacher, bool) {
	list, ok := db.GetTeachersDB()
	if !ok {
		return list, ok
	}
	for _, te := range list {
		db.FillTeacherCareersFromLegacy(te)
	}
	return list, true
}

/* InsertTeacherService newPortalPwd: contraseña de acceso al portal (opcional; mín. 6 caracteres si se informa). */
func InsertTeacherService(t models.Teacher, newPortalPwd string) (string, int, error) {
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
	if msg, ok := validateTeacherCareers(&t); !ok {
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
		t.Password = h
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

/* UpdateTeacherService newPortalPwd: si no está vacío, actualiza contraseña de portal del docente. */
func UpdateTeacherService(t models.Teacher, newPortalPwd string) (string, int, error) {
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
	if msg, ok := validateTeacherCareers(&t); !ok {
		return msg, 199, nil
	}

	t.UpdatedAt = time.Now()
	t.Password = ""
	okDB, err := db.UpdateTeacherDB(t)
	if err != nil || !okDB {
		return "No se pudo actualizar el docente", 400, err
	}
	if np := strings.TrimSpace(newPortalPwd); np != "" {
		if len(np) < 6 {
			return "La contraseña de acceso al portal debe tener al menos 6 caracteres", 199, nil
		}
		h, errH := utils.EncryptPassword(np)
		if errH != nil {
			return "Error al cifrar la contraseña", 400, errH
		}
		if errP := db.UpdateTeacherPasswordHashDB(t.ID, h); errP != nil {
			return "Docente actualizado pero no se pudo guardar la contraseña", 400, errP
		}
	}
	return "El docente se actualizo correctamente", 200, nil
}

const msgTeacherTieneAsignaciones = "el docente tiene asignaciones, por favor limpiar las mismas"

/* UpdateTeacherActiveService solo activo/inactivo; no permite pasar a inactivo si tiene carreras asignadas */
func UpdateTeacherActiveService(id primitive.ObjectID, active bool) (string, int, error) {
	if id.IsZero() {
		return "Falta el id del docente", 199, nil
	}
	t, ok := db.GetTeacherByID(id)
	if !ok {
		return "Docente no encontrado", 199, nil
	}
	if !active {
		db.FillTeacherCareersFromLegacy(&t)
		if len(t.Careers) > 0 {
			return msgTeacherTieneAsignaciones, 199, nil
		}
	}
	okDB, err := db.UpdateTeacherActiveDB(id, active)
	if err != nil || !okDB {
		return "No se pudo actualizar el estado del docente", 400, err
	}
	return "El estado del docente se actualizo correctamente", 200, nil
}

