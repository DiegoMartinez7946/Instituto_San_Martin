package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* StudentRoutes */
func StudentRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/student", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetStudents))).Methods("GET")
	router.HandleFunc("/student", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertStudent))).Methods("POST")
	router.HandleFunc("/student", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.UpdateStudent))).Methods("PUT")
	router.HandleFunc("/student/active", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.ChangeActiveStudent))).Methods("PUT")
	return router
}
