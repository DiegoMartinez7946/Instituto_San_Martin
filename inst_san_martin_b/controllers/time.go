package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/services"
)

/* GetArgentinaClock GET /time/argentina — hora oficial Argentina (IANA), persistida en colección time */
func GetArgentinaClock(w http.ResponseWriter, r *http.Request) {
	data, _, _ := services.GetArgentinaClockService()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Response{
		Message: "Hora Argentina",
		Code:    200,
		Ok:      true,
		Data:    data,
	})
}
