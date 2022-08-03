package services

import (
	"regexp"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

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

	// verify if the name has already exists
	_, check, errorCheck := db.CheckExistDegree(d.Name)
	if check == true {
		return "Ya existe esa carrera en el sistema", 199, errorCheck
	}

	row := models.Degree {
		Name: d.Name,
		Active: d.Active,
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
	if len(d.Name) == 0 {
		return "La carrera no puede venir vacia", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, d.Name)
	if anyNumber == true {
		return "No puede actualizar la carrera con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistDegree(d.Name)
	if check == true {
		return "Ya existe esa carrera en el sistema", 199, errorCheck
	}

	_, err := db.UpdateDegreeDB(d)
	if err != nil {
		return "Hubo un error al actualizar la carrera en la base", 400, err
	}

	return "La carrera se actualizo correctamente", 200, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateDegreeStatusService update the academy degree active */
func UpdateDegreeStatusService(d models.Degree) (string, int, error) {

	_, err := db.UpdateStatusDegreeDB(d)
	if err != nil {
		return "Hubo un error al actualizar el estado de la carrera en la base", 400, err
	}

	return "El estado de la carrera se actualizo correctamente", 200, nil
}
