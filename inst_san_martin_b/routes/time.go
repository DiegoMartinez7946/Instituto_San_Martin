package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* TimeRoutes reloj de referencia Argentina */
func TimeRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/time/argentina", middlewares.DbCheck(middlewares.ValidatedJWT(controllers.GetArgentinaClock))).Methods("GET")
	return router
}
