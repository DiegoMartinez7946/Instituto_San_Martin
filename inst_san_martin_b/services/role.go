package services

import (
	"log"
	"regexp"
	"strings"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/* GetRolesService call the db to get the roles */
func GetRolesService() ([]*models.Role, bool) {
	// call the db
	result, status := db.GetRolesDB()
	if status == false {
		log.Println("There was an error in services -> GetRolesService")
		return result, status
	}

	return result, status
}

/* Invice call the db to insert the role */
func InsertRoleService(r models.Role) (string, bool, error) {
	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, r.Type)
	if anyNumber == true {
		log.Println("No puede registrar el tipo de rol con numeros")
		return "", false, errRegexp
	}

	// verify if the type has already exists
	_, check, errorCheck := db.CheckExistRole(r.Type)
	if check == true {
		log.Println("Ya existe ese rol en el sistema")
		return "", false, errorCheck
	}

	r.Type = strings.ToUpper(r.Type)
	row := models.Role {
		Type: r.Type,
	}

	_, status, err := db.InsertRoleDB(row)
	if status == false || err != nil {
		log.Println("There was an error in services -> InsertRoleDB")
		return "", false, err
	}

	return "", status, err
}

/* UpdateRoleService update the user role */
func UpdateRoleService(r models.Role) (bool, error) {
	if len(r.Type) == 0 {
		log.Println("El rol no puede venir vacio")
		return false, nil
	}

	r.Type = strings.ToUpper(r.Type)

	_, err := db.UpdateRoleDB(r)
	if err != nil {
		return false, err
	}

	return true, nil
}

/* DeleteRoleService delete the user role */
func DeleteRoleService(IDRole string) error {
	// aca despues tengo que agregar una verificacion si existe un usuario
	// con un rol asociado, en caso de existir no me dejaria borrarlo
	err := db.DeleteRoleDB(IDRole)
	if err != nil {
		log.Println("Hubo un error en db -> Role -> DeleteRoleDB")
	}
	return err
}
