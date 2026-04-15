package services

import (
	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/* GetModalidadesService */
func GetModalidadesService() ([]*models.ModalidadDoc, bool) {
	return db.GetModalidadesDB()
}
