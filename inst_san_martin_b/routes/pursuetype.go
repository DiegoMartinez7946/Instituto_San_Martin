// Package routes provides ...
package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* PursueType Routes */
func PursueTypeRoutes(router *mux.Router) (*mux.Router) {
	router.HandleFunc("/pursuetype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetPursueTypes))).Methods("GET")	
	router.HandleFunc("/pursuetype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertPursueType))).Methods("POST")	
	router.HandleFunc("/pursuetype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.UpdatePursueType))).Methods("PUT")	
	router.HandleFunc("/pursuetype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.DeletePursueType))).Methods("DELETE")	

	return router
}

