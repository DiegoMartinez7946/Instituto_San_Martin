// Package services provides ...
package services

import (
	"log"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
)

/* InsertUserService check the user income andthen insert in the db */
func InsertUserService(u models.User) (string, bool, error) {
	log.Println("services")
	// 1 verify existing
	user, status, err := db.CheckExistUser(u.Email)
	if user.Email != "" {
		return "Ya existe un usuario registrado con ese email!! ", status, err
	}
	// 2 check password large, correct email
	var idUserType string
	idUserType, status, err = db.CheckExistRole(u.UserType)
	if idUserType == "" {
		return "El tipo de usuario no existe!! ", status, err
	}

	u.Password, _ = utils.EncryptPassword(u.Password)
	row := models.User {
		Email: u.Email,
		Password: u.Password,
		UserType: idUserType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, status, err = db.InsertUserDB(row)
	if status == false || err != nil {
		log.Println("There was an error in services -> InsertUserDB")
		return "", status, err
	}

	return "", status, err
}

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

/* GetUsersService call the db to get the users */
func GetUsersService() ([]models.UserResponse, bool) {
	// define roleType
	roleType := "ADMINISTRATIVO"
	if GUserType == "ADMINISTRATIVO" {
		roleType = "ALUMNO"
	}
	if GUserType == "ALUMNO" {
		roleType = ""
	}
	// call the db
	result, status, error := db.GetUsersDB(roleType)
	if status == false {
		log.Println("There was an error in services -> GetUsersService " + error.Error())
		return result, status
	}

	return result, status
}

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
		Data: cp,
	}
	return resp, true, nil
}
