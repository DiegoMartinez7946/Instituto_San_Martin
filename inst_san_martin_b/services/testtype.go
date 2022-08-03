package services

import (
	"regexp"
	"strings"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/***************************************************************/
/***************************************************************/
/* GetTestTypesService call the db to get the test types */
func GetTestTypesService() ([]*models.TestType, bool) {
	// call the db
	result, status := db.GetTestTypesDB()
	if status == false {
		return result, status
	}

	return result, status
}

/***************************************************************/
/***************************************************************/
/* InsertTestTypeService call the db to insert the test type */
func InsertTestTypeService(t models.TestType) (string, int, error) {
	// check if the testtype is empty
	if len(t.Type) == 0 {
		return "No puede registrar el tipo de examen vacio", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, t.Type)
	if anyNumber == true {
		return "No puede registrar el tipo de examen con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistTestType(t.Type)
	if check == true {
		return "Ya existe ese tipo de examen en el sistema", 199, errorCheck
	}

	t.Type = strings.ToUpper(t.Type)
	row := models.TestType {
		Type: t.Type,
	}

	msg, err := db.InsertTestTypeDB(row)
	if err != nil {
		return msg, 400, err
	}

	return msg, 201, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateTestTypeService update the test type */
func UpdateTestTypeService(t models.TestType) (string, int, error) {
	if len(t.Type) == 0 {
		return "El tipo de examen no puede venir vacio", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, t.Type)
	if anyNumber == true {
		return "No puede actualizar el tipo de examen con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistTestType(t.Type)
	if check == true {
		return "Ya existe ese tipo de examen en el sistema", 199, errorCheck
	}

	t.Type = strings.ToUpper(t.Type)

	_, err := db.UpdateTestTypeDB(t)
	if err != nil {
		return "Hubo un error al actualizar el tipo de examen en la base", 400, err
	}

	return "El tipo de examen se actualizo correctamente", 200, nil
}

/***************************************************************/
/***************************************************************/
/* DeleteTestTypeService delete the test type */
func DeleteTestTypeService(IDTestType string) (string, int, error) {
	// find the test type
	_, errFind := db.GetTestTypeDB(IDTestType)
	if errFind != nil {
		return "Hubo un error al intentar localizar el tipo de examen a borrar en la base", 400, errFind
	}

	//// find if exists users with that role
	//users, code, errUser := db.GetUsersDB(role.Type)
	//if errUser != nil {
	//	return "Hubo un error al buscar los usuarios asociados al rol", code, errUser
	//}
	//if len(users) > 0 {
	//	return "No se puede borrar el rol, hay usuarios asociados al mismo", 199, nil
	//}

	err := db.DeleteTestTypeDB(IDTestType)
	if err != nil {
		return "Hubo un error al intentar borrar el tipo de examen en la base", 400, err
	}

	return "El tipo de examen se borro correctamente", 200, nil
}
