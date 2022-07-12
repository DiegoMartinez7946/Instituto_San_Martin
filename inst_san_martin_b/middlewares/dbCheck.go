package middlewares

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/config"
)

/* Check if the database is alive */
func DbCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.CheckConnection() == 0 {
			// si falla la conexion chau, no llega al router
			http.Error(w, "Conexion perdida con la base de datos", 500)
			return
		}
		next.ServeHTTP(w, r) // esto continua el middleware
	}
}
