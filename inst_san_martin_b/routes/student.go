package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* StudentRoutes — solo rol ADMINISTRATIVO */
func StudentRoutes(router *mux.Router) *mux.Router {
	chain := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.SoloAdministrativo(h)))
	}
	router.HandleFunc("/student", chain(controllers.GetStudents)).Methods("GET")
	router.HandleFunc("/student", chain(controllers.InsertStudent)).Methods("POST")
	router.HandleFunc("/student", chain(controllers.UpdateStudent)).Methods("PUT")
	router.HandleFunc("/student/password", chain(controllers.ResetStudentPassword)).Methods("PUT")
	router.HandleFunc("/student/active", chain(controllers.ChangeActiveStudent)).Methods("PUT")
	return router
}
