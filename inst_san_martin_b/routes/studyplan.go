package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* StudyPlanRoutes — lectura administrador o administrativo; alta/edición solo administrativo */
func StudyPlanRoutes(router *mux.Router) *mux.Router {
	chain := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.SoloAdministrativo(h)))
	}
	chainRead := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministrativoOnly(h)))
	}
	router.HandleFunc("/studyplan", chainRead(controllers.GetStudyPlans)).Methods("GET")
	router.HandleFunc("/studyplan", chain(controllers.InsertStudyPlan)).Methods("POST")
	router.HandleFunc("/studyplan", chain(controllers.UpdateStudyPlan)).Methods("PUT")
	router.HandleFunc("/studyplan/active", chain(controllers.ChangeActiveStudyPlan)).Methods("PUT")
	return router
}
