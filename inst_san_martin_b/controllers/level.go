package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* GetLevels GET /level — filas nivel + jerarquia (colección level) */
func GetLevels(w http.ResponseWriter, r *http.Request) {
	result, ok := services.GetLevelsService()
	if !ok {
		res := models.Response{
			Message: "Error al consultar los niveles",
			Code:    400,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	res := models.Response{Data: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
