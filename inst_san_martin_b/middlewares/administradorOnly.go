package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* AdministradorOnly permite la petición solo si el JWT corresponde a rol ADMINISTRADOR (gestión de cuentas user). */
func AdministradorOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		role := strings.TrimSpace(strings.ToUpper(services.GUserType))
		if role != "ADMINISTRADOR" {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(models.Response{
				Message: "Solo el rol administrador puede administrar el módulo de usuarios del sistema",
				Code:    403,
			})
			return
		}
		next.ServeHTTP(w, r)
	}
}
