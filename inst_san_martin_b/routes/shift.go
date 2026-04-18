// Package routes provides ...
package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* ShiftRoutes — lectura admin o administrativo (asignación en usuarios); ABM solo ADMINISTRATIVO */
func ShiftRoutes(router *mux.Router) *mux.Router {
	read := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministrativoOnly(h)))
	}
	write := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.SoloAdministrativo(h)))
	}
	router.HandleFunc("/shift", read(controllers.GetShifts)).Methods("GET")
	router.HandleFunc("/shift", write(controllers.InsertShift)).Methods("POST")
	router.HandleFunc("/shift", write(controllers.UpdateShift)).Methods("PUT")
	router.HandleFunc("/shift", write(controllers.DeleteShift)).Methods("DELETE")
	return router
}

