package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* TeacherRoutes — solo rol ADMINISTRATIVO (gestión docente) */
func TeacherRoutes(router *mux.Router) *mux.Router {
	chain := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministrativoOnly(h)))
	}

	router.HandleFunc("/teacher", chain(controllers.GetTeachers)).Methods("GET")
	router.HandleFunc("/teacher", chain(controllers.InsertTeacher)).Methods("POST")
	router.HandleFunc("/teacher", chain(controllers.UpdateTeacher)).Methods("PUT")

	router.HandleFunc("/titulo-habilitante", chain(controllers.GetTitulosHabilitantes)).Methods("GET")
	router.HandleFunc("/modalidad-docente", chain(controllers.GetModalidadesDoc)).Methods("GET")

	return router
}
