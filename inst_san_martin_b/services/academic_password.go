package services

import (
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* ResetTeacherPasswordAdministrativo asigna nueva contraseña de portal (solo rol ADMINISTRATIVO vía middleware). */
func ResetTeacherPasswordAdministrativo(idHex, newPassword string) (string, int, error) {
	np := strings.TrimSpace(newPassword)
	if len(np) < 6 {
		return "La contraseña debe tener al menos 6 caracteres", 199, nil
	}
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idHex))
	if err != nil || id.IsZero() {
		return "Id de docente invalido", 400, err
	}
	_, ok := db.GetTeacherByID(id)
	if !ok {
		return "Docente no encontrado", 404, nil
	}
	h, err := utils.EncryptPassword(np)
	if err != nil {
		return "Error al cifrar la contraseña", 400, err
	}
	if err := db.UpdateTeacherPasswordHashDB(id, h); err != nil {
		return "No se pudo actualizar la contraseña", 400, err
	}
	return "Contraseña del docente actualizada", 200, nil
}

/* ResetStudentPasswordAdministrativo asigna nueva contraseña de portal (solo rol ADMINISTRATIVO vía middleware). */
func ResetStudentPasswordAdministrativo(idHex, newPassword string) (string, int, error) {
	np := strings.TrimSpace(newPassword)
	if len(np) < 6 {
		return "La contraseña debe tener al menos 6 caracteres", 199, nil
	}
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idHex))
	if err != nil || id.IsZero() {
		return "Id de alumno invalido", 400, err
	}
	_, ok := db.FindStudentByIDDB(id)
	if !ok {
		return "Alumno no encontrado", 404, nil
	}
	h, err := utils.EncryptPassword(np)
	if err != nil {
		return "Error al cifrar la contraseña", 400, err
	}
	if err := db.UpdateStudentPasswordHashDB(id, h); err != nil {
		return "No se pudo actualizar la contraseña", 400, err
	}
	return "Contraseña del alumno actualizada", 200, nil
}

func changePasswordTeacherSelf(cp models.OldNewPassword) (models.Response, bool, error) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(GUserID))
	if err != nil || id.IsZero() {
		resp := models.Response{Message: "Sesion invalida", Code: 400, Ok: false}
		return resp, false, nil
	}
	t, ok := db.GetTeacherByID(id)
	if !ok || !strings.EqualFold(strings.TrimSpace(t.Email), strings.TrimSpace(cp.Email)) {
		resp := models.Response{Message: "Docente no encontrado", Code: 404, Ok: false}
		return resp, false, nil
	}
	if strings.TrimSpace(t.Password) == "" {
		resp := models.Response{Message: "No hay contraseña asignada. Solicite al administrativo que le asigne una.", Code: 400, Ok: false}
		return resp, false, nil
	}
	if err := utils.DecryptPassword(t.Password, cp.CurrentPassword); err != nil {
		resp := models.Response{Message: "Contraseña actual incorrecta", Code: 400, Ok: false}
		return resp, false, err
	}
	if len(strings.TrimSpace(cp.NewPassword)) < 6 {
		resp := models.Response{Message: "La nueva contraseña debe tener al menos 6 caracteres", Code: 400, Ok: false}
		return resp, false, nil
	}
	h, err := utils.EncryptPassword(cp.NewPassword)
	if err != nil {
		resp := models.Response{Message: "Error al cifrar la nueva contraseña", Code: 400, Ok: false}
		return resp, false, err
	}
	if err := db.UpdateTeacherPasswordHashDB(id, h); err != nil {
		resp := models.Response{Message: "No se pudo guardar la contraseña", Code: 400, Ok: false}
		return resp, false, err
	}
	return models.Response{Message: "Password actualizada", Code: 200, Ok: true}, true, nil
}

func changePasswordStudentSelf(cp models.OldNewPassword) (models.Response, bool, error) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(GUserID))
	if err != nil || id.IsZero() {
		resp := models.Response{Message: "Sesion invalida", Code: 400, Ok: false}
		return resp, false, nil
	}
	st, ok := db.FindStudentByIDDB(id)
	if !ok || !strings.EqualFold(strings.TrimSpace(st.Email), strings.TrimSpace(cp.Email)) {
		resp := models.Response{Message: "Alumno no encontrado", Code: 404, Ok: false}
		return resp, false, nil
	}
	if strings.TrimSpace(st.Password) == "" {
		resp := models.Response{Message: "No hay contraseña asignada. Solicite al administrativo que le asigne una.", Code: 400, Ok: false}
		return resp, false, nil
	}
	if err := utils.DecryptPassword(st.Password, cp.CurrentPassword); err != nil {
		resp := models.Response{Message: "Contraseña actual incorrecta", Code: 400, Ok: false}
		return resp, false, err
	}
	if len(strings.TrimSpace(cp.NewPassword)) < 6 {
		resp := models.Response{Message: "La nueva contraseña debe tener al menos 6 caracteres", Code: 400, Ok: false}
		return resp, false, nil
	}
	h, err := utils.EncryptPassword(cp.NewPassword)
	if err != nil {
		resp := models.Response{Message: "Error al cifrar la nueva contraseña", Code: 400, Ok: false}
		return resp, false, err
	}
	if err := db.UpdateStudentPasswordHashDB(id, h); err != nil {
		resp := models.Response{Message: "No se pudo guardar la contraseña", Code: 400, Ok: false}
		return resp, false, err
	}
	return models.Response{Message: "Password actualizada", Code: 200, Ok: true}, true, nil
}

func changePasswordUserStaff(cp models.OldNewPassword) (models.Response, bool, error) {
	user, status, err := db.CheckExistUser(cp.Email)
	if status == false {
		resp := models.Response{
			Message: "El usuario no esta registrado",
			Code:    404,
			Ok:      false,
		}
		return resp, false, err
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

	return models.Response{Message: "Password actualizada", Code: 200, Ok: true}, true, nil
}
