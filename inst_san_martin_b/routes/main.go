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
	router = AdminMetricsRoutes(router)
	router = TimeRoutes(router)
	router = UserRoutes(router)
	router = RoleRoutes(router)
	router = ShiftRoutes(router)
	router = PursueTypeRoutes(router)
	router = TestTypeRoutes(router)
	router = DegreeRoutes(router)
	router = LevelRoutes(router)
	router = StudyPlanRoutes(router)
	router = StudentRoutes(router)
	router = TeacherRoutes(router)

	return router
}
