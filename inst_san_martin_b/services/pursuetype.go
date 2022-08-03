package services

import (
	"regexp"
	"strings"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/***************************************************************/
/***************************************************************/
/* GetPursueTypesService call the db to get the pursue types */
func GetPursueTypesService() ([]*models.PursueType, bool) {
	// call the db
	result, status := db.GetPursueTypesDB()
	if status == false {
		return result, status
	}

	return result, status
}

/***************************************************************/
/***************************************************************/
/* InsertPursueTypeService call the db to insert the pursue type */
func InsertPursueTypeService(p models.PursueType) (string, int, error) {
	// check if the pursuetype is empty
	if len(p.Type) == 0 {
		return "No puede registrar la modalidad de cursado vacia", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, p.Type)
	if anyNumber == true {
		return "No puede registrar la modalidad de cursado con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistPursueType(p.Type)
	if check == true {
		return "Ya existe esa modalidad de cursado en el sistema", 199, errorCheck
	}

	p.Type = strings.ToUpper(p.Type)
	row := models.PursueType {
		Type: p.Type,
	}

	msg, err := db.InsertPursueTypeDB(row)
	if err != nil {
		return msg, 400, err
	}

	return msg, 201, nil
}

/***************************************************************/
/***************************************************************/
/* UpdatePursueTypeService update the pursue type */
func UpdatePursueTypeService(p models.PursueType) (string, int, error) {
	if len(p.Type) == 0 {
		return "La modalidad de cursado no puede venir vacio", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, p.Type)
	if anyNumber == true {
		return "No puede actualizar la modalidad de cursado con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistPursueType(p.Type)
	if check == true {
		return "Ya existe esa modalidad de cursado en el sistema", 199, errorCheck
	}

	p.Type = strings.ToUpper(p.Type)

	_, err := db.UpdatePursueTypeDB(p)
	if err != nil {
		return "Hubo un error al actualizar la modalidad de cursado en la base", 400, err
	}

	return "La modalidad de cursado se actualizo correctamente", 200, nil
}

/***************************************************************/
/***************************************************************/
/* DeletePursueTypeService delete the pursue type */
func DeletePursueTypeService(IDPursueType string) (string, int, error) {
	// find the pursue type
	_, errFind := db.GetPursueTypeDB(IDPursueType)
	if errFind != nil {
		return "Hubo un error al intentar localizar la modalidad de cursado a borrar en la base", 400, errFind
	}

	// find if exists users with that role
	//users, code, errUser := db.GetUsersDB(role.Type)
	//if errUser != nil {
	//	return "Hubo un error al buscar los usuarios asociados al rol", code, errUser
	//}
	//if len(users) > 0 {
	//	return "No se puede borrar el rol, hay usuarios asociados al mismo", 199, nil
	//}

	err := db.DeletePursueTypeDB(IDPursueType)
	if err != nil {
		return "Hubo un error al intentar borrar la modalidad de cursado en la base", 400, err
	}

	return "La modalidad de cursado se borro correctamente", 200, nil
}
