package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var allowedDegreeNiveles = map[string]struct{}{
	"inicial":       {},
	"primario":      {},
	"secundario":    {},
	"terciario":     {},
	"universitario": {},
}

var allowedDegreeTurnoCodes = map[string]struct{}{
	"manana": {},
	"tarde":  {},
	"noche":  {},
}

var degreeTurnoStableOrder = []string{"manana", "tarde", "noche"}

func normalizeDegreeTurnos(in []string) ([]string, string) {
	seen := make(map[string]struct{})
	for _, raw := range in {
		c := strings.ToLower(strings.TrimSpace(raw))
		if c == "" {
			continue
		}
		if _, ok := allowedDegreeTurnoCodes[c]; !ok {
			return nil, "Solo se permiten los turnos: mañana, tarde y noche"
		}
		seen[c] = struct{}{}
	}
	if len(seen) == 0 {
		return nil, "Seleccione al menos un turno en el que se dicta la carrera"
	}
	out := make([]string, 0, len(seen))
	for _, c := range degreeTurnoStableOrder {
		if _, ok := seen[c]; ok {
			out = append(out, c)
		}
	}
	return out, ""
}

func normalizeNivel(n string) string {
	return strings.ToLower(strings.TrimSpace(n))
}

func isAllowedDegreeNivel(n string) bool {
	_, ok := allowedDegreeNiveles[normalizeNivel(n)]
	return ok
}

/***************************************************************/
/***************************************************************/
/* GetDegreesService call the db to get the degrees */
func GetDegreesService() ([]*models.Degree, bool) {
	// call the db
	result, status := db.GetDegreesDB()
	if status == false {
		return result, status
	}

	return result, status
}

/***************************************************************/
/***************************************************************/
/* InsertDegreeService call the db to insert the degree */
func InsertDegreeService(d models.Degree) (string, int, error) {
	// check if the degree is empty
	if len(d.Name) == 0 {
		return "No puede registrar la carrera vacia", 199, nil
	}

	// check if the degree is active
	if d.Active != true {
		return "La carrera debe estar activa al crearse", 199, nil
	}

	// verify if the name has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, d.Name)
	if anyNumber == true {
		return "No puede registrar la carrera con numeros", 199, errRegexp
	}

	if db.IsDegreeNameTakenByOther(d.Name, primitive.ObjectID{}) {
		return "Ya existe esa carrera en el sistema", 199, nil
	}

	nivel := normalizeNivel(d.Nivel)
	if !isAllowedDegreeNivel(nivel) {
		return "Debe seleccionar un nivel valido para la carrera", 199, nil
	}
	turnos, msgTurnos := normalizeDegreeTurnos(d.Turnos)
	if msgTurnos != "" {
		return msgTurnos, 199, nil
	}
	studyNR := strings.TrimSpace(d.StudyPlanID)
	if studyNR != "" && !db.StudyPlanActiveByNumeroResolucion(studyNR) {
		return "El plan de estudio indicado no existe o no esta activo", 199, nil
	}

	row := models.Degree{
		Name:        d.Name,
		Nivel:       nivel,
		StudyPlanID: studyNR,
		Turnos:      turnos,
		Active:      d.Active,
	}

	msg, err := db.InsertDegreeDB(row)
	if err != nil {
		return msg, 400, err
	}

	return msg, 201, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateDegreeService update the academy degree */
func UpdateDegreeService(d models.Degree) (string, int, error) {
	if d.ID.IsZero() {
		return "Id de carrera invalido o no enviado", 199, nil
	}
	if len(d.Name) == 0 {
		return "La carrera no puede venir vacia", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, d.Name)
	if anyNumber == true {
		return "No puede actualizar la carrera con numeros", 199, errRegexp
	}

	if db.IsDegreeNameTakenByOther(d.Name, d.ID) {
		return "Ya existe esa carrera en el sistema", 199, nil
	}

	nivel := normalizeNivel(d.Nivel)
	if !isAllowedDegreeNivel(nivel) {
		return "Debe seleccionar un nivel valido para la carrera", 199, nil
	}
	turnos, msgTurnos := normalizeDegreeTurnos(d.Turnos)
	if msgTurnos != "" {
		return msgTurnos, 199, nil
	}
	studyNR := strings.TrimSpace(d.StudyPlanID)
	if studyNR != "" && !db.StudyPlanActiveByNumeroResolucion(studyNR) {
		return "El plan de estudio indicado no existe o no esta activo", 199, nil
	}
	d.Nivel = nivel
	d.StudyPlanID = studyNR
	d.Turnos = turnos

	_, err := db.UpdateDegreeDB(d)
	if err != nil {
		if errors.Is(err, db.ErrDegreeNotFound) {
			return "No se encontro la carrera con el id indicado", 199, nil
		}
		return "Hubo un error al actualizar la carrera en la base", 400, err
	}

	return "La carrera se actualizo correctamente", 200, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateDegreeStatusService update the academy degree active */
func UpdateDegreeStatusService(d models.Degree) (string, int, error) {
	if d.ID.IsZero() {
		return "Id de carrera invalido o no enviado", 199, nil
	}

	_, err := db.UpdateStatusDegreeDB(d)
	if err != nil {
		if errors.Is(err, db.ErrDegreeNotFound) {
			return "No se encontro la carrera con el id indicado", 199, nil
		}
		return "Hubo un error al actualizar el estado de la carrera en la base", 400, err
	}

	return "El estado de la carrera se actualizo correctamente", 200, nil
}
