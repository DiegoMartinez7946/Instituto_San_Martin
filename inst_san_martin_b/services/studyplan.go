package services

import (
	"regexp"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/***************************************************************/
/***************************************************************/
/* GetStudyPlansService call the db to get the study plans */
func GetStudyPlansService() ([]*models.StudyPlan, bool) {
	// call the db
	result, status := db.GetStudyPlansDB()
	if status == false {
		return result, status
	}

	return result, status
}

/***************************************************************/
/***************************************************************/
/* InsertStudyPlanService call the db to insert the study plan */
func InsertStudyPlanService(s models.StudyPlan) (string, int, error) {
	// check if the study plan is empty
	if len(s.Name) == 0 {
		return "No puede registrar un nombre de plan de estudio vacio", 199, nil
	}

	// check if the code of study plan is empty
	if len(s.Code) == 0 {
		return "No puede registrar un codigo de plan de estudio vacio", 199, nil
	}

	// check if the code of study plan is empty
	if len(s.DegreeId) == 0 {
		return "No puede registrar plan de estudio sin una carrera asociada", 199, nil
	}

	// check if the degree is active
	if s.State != true {
		return "El plan de estudio debe estar activo al crearse", 199, nil
	}

	// verify if the name has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, s.Name)
	if anyNumber == true {
		return "No puede registrar el nombre del plan de estudio con numeros", 199, errRegexp
	}

	// verify if the code has already exists
	_, check, errorCheck := db.CheckExistStudyPlan(s.Code)
	if check == true {
		return "Ya existe ese codigo de plan de estudio en el sistema", 199, errorCheck
	}

	row := models.StudyPlan {
		Name: s.Name,
		Code: s.Code,
		DegreeId: s.DegreeId,
		State: s.State,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	msg, err := db.InsertStudyPlanDB(row)
	if err != nil {
		return msg, 400, err
	}

	return msg, 201, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateStudyPlanService update the study plan */
func UpdateStudyPlanService(s models.StudyPlan) (string, int, error) {
	if len(s.Name) == 0 {
		return "El plan de estudio no puede venir vacio", 199, nil
	}

	// check if the code of study plan is empty
	if len(s.Code) == 0 {
		return "El codigo de plan de estudio no puede venir vacio", 199, nil
	}

	// check if the code of study plan is empty
	if len(s.DegreeId) == 0 {
		return "Debe venir una carrera al plan de estudio asociada", 199, nil
	}
	// verify if the type has any number
	anyNumber, errRegexp := regexp.MatchString(`\d+`, s.Name)
	if anyNumber == true {
		return "No puede actualizar el plan de estudios con numeros", 199, errRegexp
	}

	// verify if the code has already exists
	_, check, errorCheck := db.CheckExistStudyPlan(s.Code)
	if check == true {
		return "Ya existe ese codigo de plan de estudio en el sistema", 199, errorCheck
	}

	_, err := db.UpdateStudyPlanDB(s)
	if err != nil {
		return "Hubo un error al actualizar el plan de estudio en la base", 400, err
	}

	return "El plan de estudio se actualizo correctamente", 200, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateStudyPlanStateService update the study plan state */
func UpdateStudyPlanStateService(s models.StudyPlan) (string, int, error) {

	_, err := db.UpdateStateStudyPlanDB(s)
	if err != nil {
		return "Hubo un error al actualizar el estado del plan de estudio en la base", 400, err
	}

	return "El estado del plan de estudio se actualizo correctamente", 200, nil
}
