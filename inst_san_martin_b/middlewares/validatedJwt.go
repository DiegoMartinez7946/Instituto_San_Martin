package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* ValidatedJWT exige JWT válido y sesión resuelta (user, docente o alumno). */
func ValidatedJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, find, _, err := services.ProcessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Error en el token ! "+err.Error(), http.StatusBadRequest)
			return
		}
		if !find {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(models.Response{
				Message: "Token invalido o sesion revocada",
				Code:    401,
			})
			return
		}
		next.ServeHTTP(w, r)
	}
}
