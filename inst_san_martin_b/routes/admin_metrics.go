package routes

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/benjacifre10/san_martin_b/middlewares"
	"github.com/gorilla/mux"
)

/* AdminMetricsRoutes métricas MongoDB solo rol ADMINISTRADOR */
func AdminMetricsRoutes(router *mux.Router) *mux.Router {
	admin := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.DbCheck(middlewares.ValidatedJWT(middlewares.AdministradorOnly(h)))
	}
	router.HandleFunc("/admin/db-metrics", admin(controllers.GetDBMetrics)).Methods("GET")
	router.HandleFunc("/admin/system-metrics", admin(controllers.GetSystemMetrics)).Methods("GET")
	return router
}
