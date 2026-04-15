package services

import (
	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/* GetTitulosHabilitantesService */
func GetTitulosHabilitantesService() ([]*models.TituloHabilitante, bool) {
	return db.GetTitulosHabilitantesDB()
}
