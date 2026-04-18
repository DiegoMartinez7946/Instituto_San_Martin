package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* GetSystemMetrics GET /admin/system-metrics — solo ADMINISTRADOR (middleware). */
func GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	payload, msg, code := services.GetSystemMetricsService()
	res := models.Response{
		Data: payload,
		Code: code,
	}
	if msg != "" {
		res.Message = msg
	}
	if code >= 400 {
		res.Ok = false
	} else {
		res.Ok = true
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(res)
}
