// Package routes
package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* User Routes */
func UserRoutes(router *mux.Router) *mux.Router {

	adminUser := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministradorOnly(h)))
	}

	router.HandleFunc("/user", adminUser(controllers.InsertUser)).Methods("POST")
	router.HandleFunc("/login", middlewares.DbCheck(controllers.Login)).Methods("POST")
	router.HandleFunc("/user", adminUser(controllers.GetUsers)).Methods("GET")
	router.HandleFunc("/user", adminUser(controllers.UpdateUser)).Methods("PUT")
	router.HandleFunc("/user/active", adminUser(controllers.ChangeActiveUser)).Methods("PUT")
	router.HandleFunc("/user/password", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.ChangePassword))).Methods("PUT")
	router.HandleFunc("/user/password/blank", adminUser(controllers.BlankPassword)).Methods("PUT")

	return router
}
