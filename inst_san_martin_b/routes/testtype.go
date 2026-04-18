// Package routes provides ...
package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* TestTypeRoutes — lectura admin o administrativo; ABM solo ADMINISTRATIVO */
func TestTypeRoutes(router *mux.Router) *mux.Router {
	read := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministrativoOnly(h)))
	}
	write := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.SoloAdministrativo(h)))
	}
	router.HandleFunc("/testtype", read(controllers.GetTestTypes)).Methods("GET")
	router.HandleFunc("/testtype", write(controllers.InsertTestType)).Methods("POST")
	router.HandleFunc("/testtype", write(controllers.UpdateTestType)).Methods("PUT")
	router.HandleFunc("/testtype", write(controllers.DeleteTestType)).Methods("DELETE")
	return router
}
