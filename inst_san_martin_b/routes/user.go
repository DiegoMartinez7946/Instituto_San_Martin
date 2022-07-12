// Package routes
package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* User Routes */
func UserRoutes(router *mux.Router) (*mux.Router) {
	
	router.HandleFunc("/user", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertUser))).Methods("POST")
	router.HandleFunc("/login", middlewares.DbCheck(controllers.Login)).Methods("POST")
	router.HandleFunc("/user", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetUsers))).Methods("GET")
	router.HandleFunc("/user/password", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.ChangePassword))).Methods("PUT")

	return router
}
