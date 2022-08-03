// Package services provides ...
package services

import (
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
)

/***************************************************************/
/***************************************************************/
/* InsertUserService check the user income andthen insert in the db */
func InsertUserService(u models.User) (string, int, error) {
	// 1 verify existing
	user, status, err := db.CheckExistUser(u.Email)
	if user.Email != "" {
		return "Ya existe un usuario registrado con ese email!! ", 400, err
	}
	// 2 check password large, correct email
	var idUserType string
	idUserType, status, err = db.CheckExistRole(u.UserType)
	if idUserType == "" {
		return "El tipo de usuario no existe!! ", 400, err
	}

	u.Password, _ = utils.EncryptPassword(u.Password)
	row := models.User {
		Email: u.Email,
		Password: u.Password,
		UserType: idUserType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	msg, err := db.InsertUserDB(row)
	if status == false || err != nil {
		return msg, 400, err
	}

	return msg, 201, err
}

/***************************************************************/
/***************************************************************/
/* LoginService check the user and create the token */
func LoginService(u models.User) (models.LoginResponse, bool, error) {
	var resp models.LoginResponse
	// find the user
	user, status, err := db.CheckExistUser(u.Email)
	if status == false {
		return resp, status, err 
	}
	// check the password
	errPassword := utils.DecryptPassword(user.Password, u.Password)
	if errPassword != nil {
		return resp, false, errPassword
	}

	// get the user with role
	userClaim, status, err := db.GetUserDB(u.Email)
	if status == false {
		return resp, status, err
	}

	// generate the token
	jwtKey, err := utils.GenerateJWT(userClaim)
	if err != nil {
		return resp, false, err		
 	}

	resp = models.LoginResponse {
		Token: jwtKey,
	}
	return resp, true, nil
}

/***************************************************************/
/***************************************************************/
/* GetUsersService call the db to get the users */
func GetUsersService() ([]models.UserResponse, bool, error) {
	// define roleType
	roleType := "ADMINISTRATIVO"
	if GUserType == "ADMINISTRATIVO" {
		roleType = "ALUMNO"
	}
	if GUserType == "ALUMNO" {
		roleType = ""
	}
	// call the db
	result, code, err := db.GetUsersDB(roleType)
	if code == 400 {
		return result, false, err
	}

	return result, true, nil
}

/***************************************************************/
/***************************************************************/
/* ChangePasswordService check the current password an update with the new one */
func ChangePasswordService(cp models.OldNewPassword) (models.Response, bool, error) {
	// find the user
	user, status, err := db.CheckExistUser(cp.Email)
	if status == false {
		resp := models.Response {
			Message: "El usuario no esta registrado",
			Code: 404,
			Ok: false,
		}
		return resp, status, err 
	}
	// check the password
	errPassword := utils.DecryptPassword(user.Password, cp.CurrentPassword)
	if errPassword != nil {
		resp := models.Response {
			Message: "Error: " + errPassword.Error(),
			Code: 404,
			Ok: false,
		}
		return resp, false, errPassword
	}

	cp.NewPassword, _ = utils.EncryptPassword(cp.NewPassword)
	row := models.User {
		Email: cp.Email,
		Password: cp.NewPassword,
		UpdatedAt: time.Now(),
	}
	
	status, err = db.ChangePasswordDB(row);
	if err != nil {
		resp := models.Response {
			Message: "Error: " + err.Error(),
			Code: 404,
			Ok: false,
		}
		return resp, false, errPassword
	}

	resp := models.Response {
		Message: "Password actualizada",
		Code: 200,
		Ok: true,
	}
	return resp, true, nil
}

/***************************************************************/
/***************************************************************/
func BlankPasswordServices(cp models.OldNewPassword) (string, int, error) {
	if len(cp.Email) == 0 {
		return "El email no puede venir vacio", 199, nil
	}

	if len(cp.NewPassword) < 6 {
		return "El password no puede tener menos de 6 caracteres", 199, nil
	}

	// find the user
	_, status, err := db.CheckExistUser(cp.Email)
	if status == false {
		return "El usuario no esta registrado", 404, err 
	}

	cp.NewPassword, _ = utils.EncryptPassword(cp.NewPassword)
	row := models.User {
		Email: cp.Email,
		Password: cp.NewPassword,
		UpdatedAt: time.Now(),
	}

	status, err = db.ChangePasswordDB(row);
	if err != nil {
		return "Error al actualizar la password", 404, err
	}

	return "Password actualizada", 200, err
}
