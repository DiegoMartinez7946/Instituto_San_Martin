// Package routes provides ...
package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* TestType Routes */
func TestTypeRoutes(router *mux.Router) (*mux.Router) {
	router.HandleFunc("/testtype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetTestTypes))).Methods("GET")	
	router.HandleFunc("/testtype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertTestType))).Methods("POST")	
	router.HandleFunc("/testtype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.UpdateTestType))).Methods("PUT")	
	router.HandleFunc("/testtype", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.DeleteTestType))).Methods("DELETE")	

	return router
}
