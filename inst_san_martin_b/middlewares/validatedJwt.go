package middlewares

import (
	"net/http"

	"github.com/benjacifre10/san_martin_b/services"
)

/* Validated the incoming jwt token */
func ValidatedJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _, _, err := services.ProcessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Error en el token ! " + err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}
}
