package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* StudyPlan Routes */
func StudyPlanRoutes(router *mux.Router) (*mux.Router) {

	router.HandleFunc("/studyplan", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetStudyPlans))).Methods("GET")
	router.HandleFunc("/studyplan", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertStudyPlan))).Methods("POST")
	router.HandleFunc("/studyplan", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.UpdateStudyPlan))).Methods("PUT")
	router.HandleFunc("/studyplan/status", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.ChangeStateStudyPlan))).Methods("PUT")

	return router
}
