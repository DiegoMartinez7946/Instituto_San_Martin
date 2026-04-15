package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* AdministrativoOnly permite la petición solo si el JWT corresponde a rol ADMINISTRATIVO (docentes y datos sensibles de gestión académica). */
func AdministrativoOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.TrimSpace(strings.ToUpper(services.GUserType)) != "ADMINISTRATIVO" {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(models.Response{
				Message: "Solo usuarios administrativos pueden acceder a este recurso",
				Code:    403,
			})
			return
		}
		next.ServeHTTP(w, r)
	}
}
