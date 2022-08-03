// Package services provides ...
package services

import (
	"regexp"
	"strings"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/***************************************************************/
/***************************************************************/
/* GetShiftsService call the db to get the shifts */
func GetShiftsService() ([]*models.Shift, int, error) {
	// call the db
	result, err := db.GetShiftsDB()
	if err != nil {
		return result, 400, err
	}

	return result, 200, nil
}

/***************************************************************/
/***************************************************************/
/* InsertShiftService call the db to insert the role */
func InsertShiftService(s models.Shift) (string, int, error) {
	// check if the shift is empty
	if len(s.Type) == 0 {
		return "No puede registrar el turno vacio", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, s.Type)
	if anyNumber == true {
		return "No puede registrar el turno con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistShift(s.Type)
	if check == true {
		return "Ya existe ese turno en el sistema", 199, errorCheck
	}

	s.Type = strings.ToUpper(s.Type)
	row := models.Shift {
		Type: s.Type,
	}

	msg, err := db.InsertShiftDB(row)
	if err != nil {
		return msg, 400, err
	}

	return msg, 201, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateShiftService update the academy shift */
func UpdateShiftService(s models.Shift) (string, int, error) {
	if len(s.Type) == 0 {
		return "El turno no puede venir vacio", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, s.Type)
	if anyNumber == true {
		return "No puede actualizar el turno con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistShift(s.Type)
	if check == true {
		return "Ya existe ese turno en el sistema", 199, errorCheck
	}

	s.Type = strings.ToUpper(s.Type)

	_, err := db.UpdateShiftDB(s)
	if err != nil {
		return "Hubo un error al actualizar el turno en la base", 400, err
	}

	return "El turno se actualizo correctamente", 200, nil
}

/***************************************************************/
/***************************************************************/
/* DeleteShiftService delete the academy shift */
func DeleteShiftService(IDShift string) (string, int, error) {
	// find the shift, despues le saco el undescore y le pongo shift para habilitar la sig funcion
	_, errFind := db.GetShiftDB(IDShift)
	if errFind != nil {
		return "Hubo un error al intentar localizar el turno a borrar en la base", 400, errFind
	}

	// find if exists subject with that shift
	//subjects, code, errSubject := db.GetSubjectsDB(shift.Type)
	//if errSubject != nil {
	//	return "Hubo un error al buscar las materias asociadas a ese turno", code, errSubject
	//}
	//if len(subjects) > 0 {
	//	return "No se puede borrar el turno, hay materias asociadas al mismo", 199, nil
	//}

	err := db.DeleteShiftDB(IDShift)
	if err != nil {
		return "Hubo un error al intentar borrar el turno en la base", 400, err
	}

	return "El turno se borro correctamente", 200, nil
}
