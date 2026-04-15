package services

import (
	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/* GetJerarquiasService returns nivel → jerarquia definitions from the jerarquia collection */
func GetJerarquiasService() ([]*models.LevelHierarchy, bool) {
	return db.GetJerarquiasDB()
}
