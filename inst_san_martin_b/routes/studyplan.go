package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* StudyPlanRoutes — solo ADMINISTRATIVO (plan de estudio no corresponde al administrador de cuentas) */
func StudyPlanRoutes(router *mux.Router) *mux.Router {
	chain := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.SoloAdministrativo(h)))
	}
	router.HandleFunc("/studyplan", chain(controllers.GetStudyPlans)).Methods("GET")
	router.HandleFunc("/studyplan", chain(controllers.InsertStudyPlan)).Methods("POST")
	router.HandleFunc("/studyplan", chain(controllers.UpdateStudyPlan)).Methods("PUT")
	router.HandleFunc("/studyplan/active", chain(controllers.ChangeActiveStudyPlan)).Methods("PUT")
	return router
}
