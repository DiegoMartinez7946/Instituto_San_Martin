package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* GetModalidadesDoc catálogo modalidad (titular / provisional / suplente) */
func GetModalidadesDoc(w http.ResponseWriter, r *http.Request) {
	result, ok := services.GetModalidadesService()
	if !ok {
		res := models.Response{Message: "Error al consultar modalidades", Code: 400}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	res := models.Response{Data: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
