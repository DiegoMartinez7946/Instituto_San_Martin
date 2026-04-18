package db

import (
	"context"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbMetricsCollection = "dbmetrics"

/* InsertDBMetricSnapshot guarda un punto para tendencias / auditoría. */
func InsertDBMetricSnapshot(s models.DBMetricSnapshot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := config.WorkingDatabase().Collection(dbMetricsCollection)
	_, err := coll.InsertOne(ctx, s)
	return err
}

/* LatestDBMetricSnapshot devuelve el último documento (si existe). */
func LatestDBMetricSnapshot() (models.DBMetricSnapshot, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	coll := config.WorkingDatabase().Collection(dbMetricsCollection)
	opts := options.FindOne().SetSort(bson.D{{Key: "createdat", Value: -1}})
	var row models.DBMetricSnapshot
	err := coll.FindOne(ctx, bson.M{}, opts).Decode(&row)
	if err != nil || row.ID.IsZero() {
		return models.DBMetricSnapshot{}, false
	}
	return row, true
}
