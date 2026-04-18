package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* DegreeRoutes — lectura admin o administrativo (formulario usuarios); ABM solo ADMINISTRATIVO */
func DegreeRoutes(router *mux.Router) *mux.Router {
	read := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministrativoOnly(h)))
	}
	write := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.SoloAdministrativo(h)))
	}
	router.HandleFunc("/degree", read(controllers.GetDegrees)).Methods("GET")
	router.HandleFunc("/degree", write(controllers.InsertDegree)).Methods("POST")
	router.HandleFunc("/degree", write(controllers.UpdateDegree)).Methods("PUT")
	router.HandleFunc("/degree/active", write(controllers.ChangeActiveDegree)).Methods("PUT")
	return router
}
