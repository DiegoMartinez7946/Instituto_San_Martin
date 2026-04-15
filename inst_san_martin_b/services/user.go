// Package services provides ...
package services

import (
	"sort"
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

	degStored, msgDeg, codeDeg := normalizeUserDegreeIDs(u.DegreeIDHexes, u.UserType)
	if codeDeg != 0 {
		return msgDeg, codeDeg, nil
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
		DegreeIDs: degStored,
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
	users, code, err := db.GetUsersDB("")
	if code == 400 {
		return users, false, err
	}

	// Docentes y alumnos viven en colecciones propias (teacher / student), no en "user".
	// Se agregan aquí para un listado único; no son editables con el mismo modal que las cuentas de login.
	seenEmail := make(map[string]struct{})
	for _, u := range users {
		if em := strings.ToLower(strings.TrimSpace(u.Email)); em != "" {
			seenEmail[em] = struct{}{}
		}
	}

	out := make([]models.UserListRow, 0, len(users)+32)
	for i := range users {
		out = append(out, users[i])
	}

	if teachers, ok := db.GetTeachersDB(); ok {
		for _, t := range teachers {
			if t == nil {
				continue
			}
			em := strings.ToLower(strings.TrimSpace(t.Email))
			if em != "" {
				if _, dup := seenEmail[em]; dup {
					continue
				}
				seenEmail[em] = struct{}{}
			}
			degTeacher := make([]string, 0, len(t.Careers))
			seenDeg := make(map[string]struct{})
			for _, c := range t.Careers {
				if c.DegreeID.IsZero() {
					continue
				}
				h := c.DegreeID.Hex()
				if _, ok := seenDeg[h]; ok {
					continue
				}
				seenDeg[h] = struct{}{}
				degTeacher = append(degTeacher, h)
			}
			if len(degTeacher) == 0 && len(t.LegacyDegreeIDs) > 0 {
				for _, d := range t.LegacyDegreeIDs {
					if d.IsZero() {
						continue
					}
					h := d.Hex()
					if _, ok := seenDeg[h]; ok {
						continue
					}
					seenDeg[h] = struct{}{}
					degTeacher = append(degTeacher, h)
				}
			}
			out = append(out, models.UserListRow{
				ID:        "teacher:" + t.ID.Hex(),
				Email:     t.Email,
				Role:      "DOCENTE",
				Name:      t.Name,
				DNI:       t.DNI,
				Address:   t.Address,
				Phone:     t.Phone,
				DegreeIDs: degTeacher,
				Active:    t.Active,
				Source:    "teacher",
			})
		}
	}

	allStudents, stOk := db.GetStudentsDB()
	if stOk {
		for _, s := range allStudents {
			if s == nil {
				continue
			}
			em := strings.ToLower(strings.TrimSpace(s.Email))
			if em != "" {
				if _, dup := seenEmail[em]; dup {
					continue
				}
				seenEmail[em] = struct{}{}
			}
			degStu := make([]string, 0, len(s.DegreeIDs))
			for _, d := range s.DegreeIDs {
				if d.IsZero() {
					continue
				}
				degStu = append(degStu, d.Hex())
			}
			out = append(out, models.UserListRow{
				ID:        "student:" + s.ID.Hex(),
				Email:     s.Email,
				Role:      "ALUMNO",
				Name:      s.Name,
				DNI:       s.DNI,
				Address:   s.Address,
				Phone:     s.Phone,
				DegreeIDs: degStu,
				Active:    s.Active,
				Modalidad: s.Modalidad,
				Condicion: s.Condicion,
				Source:    "student",
			})
		}
	}

	stByEmail := make(map[string]*models.Student)
	for _, s := range allStudents {
		if s == nil {
			continue
		}
		em := strings.ToLower(strings.TrimSpace(s.Email))
		if em == "" {
			continue
		}
		stByEmail[em] = s
	}
	mergeAlumnoUserRowFromStudent := func(row *models.UserListRow, st *models.Student) {
		row.StudentLinked = true
		if st != nil && !st.ID.IsZero() {
			row.StudentRecordID = st.ID.Hex()
		}
		row.Modalidad = st.Modalidad
		row.Condicion = st.Condicion
		if len(st.DegreeIDs) > 0 {
			degHex := make([]string, 0, len(st.DegreeIDs))
			for _, d := range st.DegreeIDs {
				if d.IsZero() {
					continue
				}
				degHex = append(degHex, d.Hex())
			}
			row.DegreeIDs = degHex
		}
	}

	for i := range out {
		if out[i].Source != "" {
			continue
		}
		if strings.ToUpper(strings.TrimSpace(out[i].Role)) != "ALUMNO" {
			continue
		}
		em := strings.ToLower(strings.TrimSpace(out[i].Email))
		if em == "" {
			continue
		}
		st, ok := stByEmail[em]
		if !ok || st == nil {
			continue
		}
		mergeAlumnoUserRowFromStudent(&out[i], st)
	}

	stByDNI := make(map[string]*models.Student)
	for _, s := range allStudents {
		if s == nil {
			continue
		}
		d := strings.TrimSpace(s.DNI)
		if d == "" {
			continue
		}
		stByDNI[d] = s
	}
	for i := range out {
		if out[i].Source != "" {
			continue
		}
		if strings.ToUpper(strings.TrimSpace(out[i].Role)) != "ALUMNO" {
			continue
		}
		if out[i].StudentLinked {
			continue
		}
		d := strings.TrimSpace(out[i].DNI)
		if d == "" {
			continue
		}
		st, ok := stByDNI[d]
		if !ok || st == nil {
			continue
		}
		mergeAlumnoUserRowFromStudent(&out[i], st)
	}

	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(strings.TrimSpace(out[i].Name)) < strings.ToLower(strings.TrimSpace(out[j].Name))
	})

	return out, true, nil
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

/* normalizeUserDegreeIDs carreras para ALUMNO/DOCENTE (como máximo una; opcional vacío) */
func normalizeUserDegreeIDs(hexes []string, roleTypeCode string) ([]primitive.ObjectID, string, int) {
	roleUpper := strings.ToUpper(strings.TrimSpace(roleTypeCode))
	if roleUpper != "ALUMNO" && roleUpper != "DOCENTE" {
		return nil, "", 0
	}
	out := make([]primitive.ObjectID, 0, 1)
	for _, raw := range hexes {
		h := strings.TrimSpace(raw)
		if h == "" {
			continue
		}
		oid, err := primitive.ObjectIDFromHex(h)
		if err != nil {
			return nil, "Id de carrera invalido", 199
		}
		if !db.DegreeDocumentExists(oid) {
			return nil, "La carrera seleccionada no es valida", 199
		}
		out = append(out, oid)
	}
	if len(out) > 1 {
		return nil, "Seleccione solo una carrera", 199
	}
	return out, "", 0
}

/* resolveStudentForAlumnoSync localiza la ficha: prioriza studentRecordId (hex) si coincide email o DNI con el formulario. */
func resolveStudentForAlumnoSync(email, dni, studentRecordIDHex string) (models.Student, bool) {
	email = strings.TrimSpace(email)
	dni = strings.TrimSpace(dni)
	if h := strings.TrimSpace(studentRecordIDHex); h != "" {
		oid, err := primitive.ObjectIDFromHex(h)
		if err == nil && !oid.IsZero() {
			st, ok := db.FindStudentByIDDB(oid)
			if ok {
				em := strings.ToLower(strings.TrimSpace(st.Email))
				uem := strings.ToLower(email)
				sdni := strings.TrimSpace(st.DNI)
				udni := dni
				if (em != "" && uem != "" && em == uem) || (sdni != "" && udni != "" && sdni == udni) {
					return st, true
				}
			}
		}
	}
	return db.FindStudentForUserSync(email, dni)
}

/* Si existe ficha student para el mismo email/DNI (o id explícito), valida y normaliza modalidad/condición alineada al usuario (activo/inactivo). */
func validateAndNormalizeAlumnoStudentForUserUpdate(userID primitive.ObjectID, email, dni string, degStored []primitive.ObjectID, modalidadIn, condicionIn, studentRecordIDHex string) (msg string, code int, studentID primitive.ObjectID, modalidadOut, condicionOut string, shouldSync bool) {
	st, found := resolveStudentForAlumnoSync(email, dni, studentRecordIDHex)
	if !found {
		return "", 0, primitive.NilObjectID, "", "", false
	}
	u, okU, errU := db.FindUserByIDDB(userID)
	if errU != nil || !okU {
		return "", 0, primitive.NilObjectID, "", "", false
	}
	s := st
	s.DegreeIDs = degStored
	s.Modalidad = modalidadIn
	s.Condicion = condicionIn
	s.Active = u.Active
	if m, okVal := validateStudentExtraModeFields(&s); !okVal {
		return m, 199, primitive.NilObjectID, "", "", false
	}
	return "", 0, st.ID, s.Modalidad, s.Condicion, true
}

/***************************************************************/
/* UpdateUserService actualiza perfil de usuario de sistema */
func UpdateUserService(idHex string, name, dni, address, phone, email, password, userTypeCode string, shiftIDs []string, degreeHexes []string, modalidadIn, condicionIn, studentRecordIDHex string) (string, int, error) {
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

	effectiveRole := strings.ToUpper(strings.TrimSpace(userTypeCode))
	if effectiveRole == "" {
		uDB, okU, errU := db.FindUserByIDDB(id)
		if errU == nil && okU && strings.TrimSpace(uDB.UserType) != "" {
			roleDoc, errR := db.GetRoleDB(strings.TrimSpace(uDB.UserType))
			if errR == nil && roleDoc.Type != "" {
				effectiveRole = strings.ToUpper(strings.TrimSpace(roleDoc.Type))
			}
		}
	}

	degStored, msgDeg, codeDeg := normalizeUserDegreeIDs(degreeHexes, effectiveRole)
	if codeDeg != 0 {
		return msgDeg, codeDeg, nil
	}
	if effectiveRole != "ALUMNO" && effectiveRole != "DOCENTE" {
		degStored = []primitive.ObjectID{}
	}
	if degStored == nil {
		degStored = []primitive.ObjectID{}
	}

	var syncStudentID primitive.ObjectID
	var syncModalidad, syncCondicion string
	doSyncStudent := false
	if effectiveRole == "ALUMNO" {
		msgSt, codeSt, sid, modOut, condOut, sync := validateAndNormalizeAlumnoStudentForUserUpdate(id, email, dni, degStored, modalidadIn, condicionIn, studentRecordIDHex)
		if codeSt != 0 {
			return msgSt, codeSt, nil
		}
		if sync {
			syncStudentID = sid
			syncModalidad = modOut
			syncCondicion = condOut
			doSyncStudent = true
		}
	}

	okDB, err := db.UpdateUserProfileDB(id, name, dni, address, phone, email, shiftStored, degStored)
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

	if doSyncStudent && !syncStudentID.IsZero() {
		_, errSync := db.UpdateStudentEnrollmentFieldsDB(syncStudentID, degStored, syncModalidad, syncCondicion)
		if errSync != nil {
			return "Usuario actualizado pero no se pudo sincronizar la ficha de alumno", 400, errSync
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
