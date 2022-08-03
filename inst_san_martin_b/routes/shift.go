// Package routes provides ...
package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* Shift Routes */
func ShiftRoutes(router *mux.Router) (*mux.Router) {
	router.HandleFunc("/shift", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetShifts))).Methods("GET")	
	router.HandleFunc("/shift", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertShift))).Methods("POST")	
	router.HandleFunc("/shift", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.UpdateShift))).Methods("PUT")	
	router.HandleFunc("/shift", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.DeleteShift))).Methods("DELETE")	

	return router
}

