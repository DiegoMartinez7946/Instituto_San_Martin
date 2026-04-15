// Package services provides ...
package services

import (
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/***************************************************************/
/* InsertUserService check the user income andthen insert in the db */
func InsertUserService(u models.User) (string, int, error) {
	u.Email = strings.TrimSpace(u.Email)
	if msg, ok := ValidateCorreoElectronicoRequired(u.Email); !ok {
		return msg, 199, nil
	}
	u.Name = strings.TrimSpace(u.Name)
	u.DNI = strings.TrimSpace(u.DNI)
	u.Address = strings.TrimSpace(u.Address)
	if len(u.Name) == 0 {
		return "El nombre completo es obligatorio", 199, nil
	}
	if msg, ok := ValidateDNIAr(u.DNI); !ok {
		return msg, 199, nil
	}
	u.Phone = NormalizeTelefonoDigits(u.Phone)
	if msg, ok := ValidateTelefono(u.Phone); !ok {
		return msg, 199, nil
	}

	_, status, err := db.CheckExistUser(u.Email)
	if err == nil && status {
		return "Ya existe un usuario registrado con ese email!! ", 400, nil
	}
	taken, err := db.FindOtherUserByDNIDB(u.DNI, primitive.NilObjectID)
	if err != nil {
		return "Error al validar DNI", 400, err
	}
	if taken {
		return "Ya existe otro usuario con ese DNI", 199, nil
	}

	idUserType, _, err := db.CheckExistRole(u.UserType)
	if idUserType == "" {
		return "El tipo de usuario no existe!! ", 400, err
	}

	mergedShifts := mergeUserShiftSlice(u.ShiftIDs, u.ShiftID)
	shiftStored, msgSh, codeSh := normalizeUserShiftIDs(mergedShifts)
	if codeSh != 0 {
		return msgSh, codeSh, nil
	}

	u.Password, _ = utils.EncryptPassword(u.Password)
	row := models.User{
		Email:     u.Email,
		Password:  u.Password,
		UserType:  idUserType,
		Name:      u.Name,
		DNI:       u.DNI,
		Address:   u.Address,
		Phone:     u.Phone,
		ShiftIDs:  shiftStored,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	msg, err := db.InsertUserDB(row)
	if err != nil {
		return msg, 400, err
	}

	return msg, 201, err
}

/***************************************************************/
/* LoginService check the user and create the token */
func LoginService(u models.User) (models.LoginResponse, bool, error) {
	var resp models.LoginResponse
	user, status, err := db.CheckExistUser(u.Email)
	if status == false {
		return resp, status, err
	}
	if !user.Active {
		return resp, false, nil
	}
	errPassword := utils.DecryptPassword(user.Password, u.Password)
	if errPassword != nil {
		return resp, false, errPassword
	}

	userClaim, status, err := db.GetUserDB(u.Email)
	if status == false {
		return resp, status, err
	}

	jwtKey, err := utils.GenerateJWT(userClaim)
	if err != nil {
		return resp, false, err
	}

	resp = models.LoginResponse{
		Token: jwtKey,
	}
	return resp, true, nil
}

/***************************************************************/
/* GetUsersService call the db to get the users */
func GetUsersService() ([]models.UserListRow, bool, error) {
	roleType := "ADMINISTRATIVO"
	if GUserType == "ADMINISTRATIVO" {
		roleType = "ALUMNO"
	}
	if GUserType == "ALUMNO" {
		roleType = ""
	}
	result, code, err := db.GetUsersDB(roleType)
	if code == 400 {
		return result, false, err
	}

	return result, true, nil
}

func mergeUserShiftSlice(ids []string, legacy string) []string {
	out := make([]string, 0, len(ids)+1)
	for _, x := range ids {
		if t := strings.TrimSpace(x); t != "" {
			out = append(out, t)
		}
	}
	if t := strings.TrimSpace(legacy); t != "" {
		out = append(out, t)
	}
	return out
}

/* normalizeUserShiftIDs valida y deduplica hex de turnos; msg y code != 0 si error */
func normalizeUserShiftIDs(ids []string) (stored []string, msg string, code int) {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(ids))
	for _, raw := range ids {
		h := strings.TrimSpace(raw)
		if h == "" {
			continue
		}
		if _, ok := seen[h]; ok {
			continue
		}
		ok, err := db.ShiftExistsByHexID(h)
		if err != nil {
			return nil, "Error al validar el turno", 400
		}
		if !ok {
			return nil, "Un turno seleccionado no es valido", 199
		}
		seen[h] = struct{}{}
		out = append(out, h)
	}
	return out, "", 0
}

/***************************************************************/
/* UpdateUserService actualiza perfil de usuario de sistema */
func UpdateUserService(idHex string, name, dni, address, phone, email, password, userTypeCode string, shiftIDs []string) (string, int, error) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idHex))
	if err != nil {
		return "Id de usuario invalido", 400, nil
	}
	_, ok, err := db.FindUserByIDDB(id)
	if err != nil || !ok {
		return "Usuario no encontrado", 404, err
	}

	email = strings.TrimSpace(email)
	if msg, okVal := ValidateCorreoElectronicoRequired(email); !okVal {
		return msg, 199, nil
	}
	name = strings.TrimSpace(name)
	dni = strings.TrimSpace(dni)
	address = strings.TrimSpace(address)
	if len(name) == 0 {
		return "El nombre completo es obligatorio", 199, nil
	}
	if msg, okVal := ValidateDNIAr(dni); !okVal {
		return msg, 199, nil
	}
	phone = NormalizeTelefonoDigits(phone)
	if msg, okVal := ValidateTelefono(phone); !okVal {
		return msg, 199, nil
	}

	takenMail, err := db.CheckExistUserByEmailOtherThanID(email, id)
	if err != nil {
		return "Error al validar email", 400, err
	}
	if takenMail {
		return "Ya existe otro usuario con ese email", 199, nil
	}

	takenDni, err := db.FindOtherUserByDNIDB(dni, id)
	if err != nil {
		return "Error al validar DNI", 400, err
	}
	if takenDni {
		return "Ya existe otro usuario con ese DNI", 199, nil
	}

	if strings.TrimSpace(userTypeCode) != "" {
		newTypeID, _, errR := db.CheckExistRole(userTypeCode)
		if newTypeID == "" {
			return "El tipo de usuario no existe", 400, errR
		}
		_, err = db.UpdateUserRoleDB(id, newTypeID)
		if err != nil {
			return "Error al actualizar el rol", 400, err
		}
	}

	shiftStored, msgSh, codeSh := normalizeUserShiftIDs(shiftIDs)
	if codeSh != 0 {
		return msgSh, codeSh, nil
	}

	okDB, err := db.UpdateUserProfileDB(id, name, dni, address, phone, email, shiftStored)
	if err != nil || !okDB {
		return "No se pudo actualizar el usuario", 400, err
	}

	if len(strings.TrimSpace(password)) >= 6 {
		hashed, _ := utils.EncryptPassword(password)
		_, err = db.UpdateUserPasswordByIDDB(id, hashed)
		if err != nil {
			return "Perfil actualizado pero fallo el cambio de password", 400, err
		}
	}

	return "El usuario se actualizo correctamente", 200, nil
}

/***************************************************************/
/* UpdateUserActiveService activo / inactivo */
func UpdateUserActiveService(id primitive.ObjectID, active bool) (string, int, error) {
	if id.IsZero() {
		return "Falta el id del usuario", 199, nil
	}
	okDB, err := db.UpdateUserActiveDB(id, active)
	if err != nil || !okDB {
		return "No se pudo actualizar el estado del usuario", 400, err
	}
	return "El estado del usuario se actualizo correctamente", 200, nil
}

/***************************************************************/
/* ChangePasswordService check the current password an update with the new one */
func ChangePasswordService(cp models.OldNewPassword) (models.Response, bool, error) {
	cp.Email = strings.TrimSpace(cp.Email)
	if msg, ok := ValidateCorreoElectronicoRequired(cp.Email); !ok {
		resp := models.Response{
			Message: msg,
			Code:    400,
			Ok:      false,
		}
		return resp, false, nil
	}
	user, status, err := db.CheckExistUser(cp.Email)
	if status == false {
		resp := models.Response{
			Message: "El usuario no esta registrado",
			Code:    404,
			Ok:      false,
		}
		return resp, status, err
	}
	errPassword := utils.DecryptPassword(user.Password, cp.CurrentPassword)
	if errPassword != nil {
		resp := models.Response{
			Message: "Error: " + errPassword.Error(),
			Code:    404,
			Ok:      false,
		}
		return resp, false, errPassword
	}

	cp.NewPassword, _ = utils.EncryptPassword(cp.NewPassword)
	row := models.User{
		Email:     strings.TrimSpace(user.Email),
		Password:  cp.NewPassword,
		UpdatedAt: time.Now(),
	}

	okPwd, errPwd := db.ChangePasswordDB(row)
	if errPwd != nil {
		resp := models.Response{
			Message: "Error al actualizar: " + errPwd.Error(),
			Code:    404,
			Ok:      false,
		}
		return resp, false, errPwd
	}
	if !okPwd {
		resp := models.Response{
			Message: "No se pudo actualizar la contraseña en la base",
			Code:    404,
			Ok:      false,
		}
		return resp, false, nil
	}

	resp := models.Response{
		Message: "Password actualizada",
		Code:    200,
		Ok:      true,
	}
	return resp, true, nil
}

/***************************************************************/
func BlankPasswordServices(cp models.OldNewPassword) (string, int, error) {
	cp.Email = strings.TrimSpace(cp.Email)
	if len(cp.Email) == 0 {
		return "El email no puede venir vacio", 199, nil
	}
	if msg, ok := ValidateCorreoElectronicoRequired(cp.Email); !ok {
		return msg, 199, nil
	}

	if len(cp.NewPassword) < 6 {
		return "El password no puede tener menos de 6 caracteres", 199, nil
	}

	existing, status, err := db.CheckExistUser(cp.Email)
	if !status {
		if err != nil {
			return "Error al buscar el usuario", 400, err
		}
		return "El usuario no esta registrado", 404, nil
	}
	emailDB := strings.TrimSpace(existing.Email)
	if emailDB == "" {
		return "El usuario no esta registrado", 404, nil
	}

	cp.NewPassword, _ = utils.EncryptPassword(cp.NewPassword)
	row := models.User{
		Email:     emailDB,
		Password:  cp.NewPassword,
		UpdatedAt: time.Now(),
	}

	okPwd, errPwd := db.ChangePasswordDB(row)
	if errPwd != nil {
		return "Error al actualizar la password", 400, errPwd
	}
	if !okPwd {
		return "No se pudo actualizar: el email no coincide con ningun registro en la base", 199, nil
	}

	return "Password actualizada", 200, nil
}
