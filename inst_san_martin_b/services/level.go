package services

import (
	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

/* GetLevelsService devuelve definiciones nivel → orden (colección level) */
func GetLevelsService() ([]*models.LevelHierarchy, bool) {
	return db.GetLevelsDB()
}
