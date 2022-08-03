package services

import (
	"regexp"
	"strings"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/***************************************************************/
/***************************************************************/
/* GetRolesService call the db to get the roles */
func GetRolesService() ([]*models.Role, bool) {
	// call the db
	result, status := db.GetRolesDB()
	if status == false {
		return result, status
	}

	return result, status
}

/***************************************************************/
/***************************************************************/
/* InsertRoleService call the db to insert the role */
func InsertRoleService(r models.Role) (string, int, error) {
	// check if the role is empty
	if len(r.Type) == 0 {
		return "No puede registrar el tipo de rol vacio", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, r.Type)
	if anyNumber == true {
		return "No puede registrar el tipo de rol con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistRole(r.Type)
	if check == true {
		return "Ya existe ese rol en el sistema", 199, errorCheck
	}

	r.Type = strings.ToUpper(r.Type)
	row := models.Role {
		Type: r.Type,
	}

	msg, err := db.InsertRoleDB(row)
	if err != nil {
		return msg, 400, err
	}

	return msg, 201, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateRoleService update the user role */
func UpdateRoleService(r models.Role) (string, int, error) {
	if len(r.Type) == 0 {
		return "El rol no puede venir vacio", 199, nil
	}

	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, r.Type)
	if anyNumber == true {
		return "No puede actualizar el tipo de rol con numeros", 199, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistRole(r.Type)
	if check == true {
		return "Ya existe ese rol en el sistema", 199, errorCheck
	}

	r.Type = strings.ToUpper(r.Type)

	_, err := db.UpdateRoleDB(r)
	if err != nil {
		return "Hubo un error al actualizar el rol en la base", 400, err
	}

	return "El rol se actualizo correctamente", 200, nil
}

/***************************************************************/
/***************************************************************/
/* DeleteRoleService delete the user role */
func DeleteRoleService(IDRole string) (string, int, error) {
	// find the role
	role, errFind := db.GetRoleDB(IDRole)
	if errFind != nil {
		return "Hubo un error al intentar localizar el rol a borrar en la base", 400, errFind
	}

	// find if exists users with that role
	users, code, errUser := db.GetUsersDB(role.Type)
	if errUser != nil {
		return "Hubo un error al buscar los usuarios asociados al rol", code, errUser
	}
	if len(users) > 0 {
		return "No se puede borrar el rol, hay usuarios asociados al mismo", 199, nil
	}

	err := db.DeleteRoleDB(IDRole)
	if err != nil {
		return "Hubo un error al intentar borrar el rol en la base", 400, err
	}

	return "El rol se borro correctamente", 200, nil
}
