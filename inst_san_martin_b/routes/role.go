package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* Role Routes */
func RoleRoutes(router *mux.Router) (*mux.Router) {

	router.HandleFunc("/user/role", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetRoles))).Methods("GET")
	router.HandleFunc("/user/role", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.InsertRole))).Methods("POST")
	router.HandleFunc("/user/role", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.UpdateRole))).Methods("PUT")
	router.HandleFunc("/user/role", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.DeleteRole))).Methods("DELETE")

	return router
}
