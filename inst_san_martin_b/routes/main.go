package routes

import (
	"github.com/benjacifre10/san_martin_b/controllers"
	"github.com/gorilla/mux"
)

/* HandlerRoutes create the differents routes for my API*/
func HandlerRoutes() (*mux.Router) {
	// create the router
	router := mux.NewRouter()

	router.HandleFunc("/health", controllers.Health).Methods("GET")
	router = UserRoutes(router)
	router = RoleRoutes(router)

	return router
}
