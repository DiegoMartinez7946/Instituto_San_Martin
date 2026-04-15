package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* GetJerarquias lists hierarchy rows (colección jerarquia) */
func GetJerarquias(w http.ResponseWriter, r *http.Request) {
	result, ok := services.GetJerarquiasService()
	if !ok {
		res := models.Response{
			Message: "Error al consultar la jerarquía de niveles",
			Code:    400,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	res := models.Response{Data: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
