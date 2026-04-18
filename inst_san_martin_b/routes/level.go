package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* LevelRoutes API de lectura para nivel → orden jerárquico */
func LevelRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/level", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetLevels))).Methods("GET")
	return router
}
