package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* JerarquiaRoutes read-only API for nivel → jerarquia */
func JerarquiaRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/jerarquia", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetJerarquias))).Methods("GET")
	return router
}
