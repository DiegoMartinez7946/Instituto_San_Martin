// Package routes provides ...
package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* PursueTypeRoutes — lectura admin o administrativo; ABM solo ADMINISTRATIVO */
func PursueTypeRoutes(router *mux.Router) *mux.Router {
	read := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministrativoOnly(h)))
	}
	write := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.SoloAdministrativo(h)))
	}
	router.HandleFunc("/pursuetype", read(controllers.GetPursueTypes)).Methods("GET")
	router.HandleFunc("/pursuetype", write(controllers.InsertPursueType)).Methods("POST")
	router.HandleFunc("/pursuetype", write(controllers.UpdatePursueType)).Methods("PUT")
	router.HandleFunc("/pursuetype", write(controllers.DeletePursueType)).Methods("DELETE")
	return router
}

