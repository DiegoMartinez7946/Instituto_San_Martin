package db

import (
	"context"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const systemMetricsCollection = "systemmetrics"

/* InsertSystemMetricSnapshot guarda un punto para tendencias / auditoría. */
func InsertSystemMetricSnapshot(s models.SystemMetricSnapshot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := config.WorkingDatabase().Collection(systemMetricsCollection)
	_, err := coll.InsertOne(ctx, s)
	return err
}

/* LatestSystemMetricSnapshot último documento si existe. */
func LatestSystemMetricSnapshot() (models.SystemMetricSnapshot, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	coll := config.WorkingDatabase().Collection(systemMetricsCollection)
	opts := options.FindOne().SetSort(bson.D{{Key: "createdat", Value: -1}})
	var row models.SystemMetricSnapshot
	err := coll.FindOne(ctx, bson.M{}, opts).Decode(&row)
	if err != nil || row.ID.IsZero() {
		return models.SystemMetricSnapshot{}, false
	}
	return row, true
}
