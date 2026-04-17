package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* SoloAdministrativo permite la petición solo si el JWT corresponde a rol ADMINISTRATIVO (gestión de alumnos y docentes). */
func SoloAdministrativo(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		role := strings.TrimSpace(strings.ToUpper(services.GUserType))
		if role != "ADMINISTRATIVO" {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(models.Response{
				Message: "Solo usuarios con rol administrativo pueden acceder a este recurso",
				Code:    403,
			})
			return
		}
		next.ServeHTTP(w, r)
	}
}
