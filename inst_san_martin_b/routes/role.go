package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* RoleRoutes — ABM de roles solo ADMINISTRADOR */
func RoleRoutes(router *mux.Router) *mux.Router {
	adminOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministradorOnly(h)))
	}
	router.HandleFunc("/user/role", adminOnly(controllers.GetRoles)).Methods("GET")
	router.HandleFunc("/user/role", adminOnly(controllers.InsertRole)).Methods("POST")
	router.HandleFunc("/user/role", adminOnly(controllers.UpdateRole)).Methods("PUT")
	router.HandleFunc("/user/role", adminOnly(controllers.DeleteRole)).Methods("DELETE")
	return router
}
