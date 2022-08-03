package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* Degree Routes */
func DegreeRoutes(router *mux.Router) (*mux.Router) {

	router.HandleFunc("/degree", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetDegrees))).Methods("GET")
	router.HandleFunc("/degree", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertDegree))).Methods("POST")
	router.HandleFunc("/degree", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.UpdateDegree))).Methods("PUT")
	router.HandleFunc("/degree/active", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.ChangeActiveDegree))).Methods("PUT")

	return router
}
