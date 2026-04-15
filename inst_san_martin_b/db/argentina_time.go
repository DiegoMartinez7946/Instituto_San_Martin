package db

import (
	"context"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeCollection = "time"

// ID fijo para un único documento de reloj (referencia Argentina).
var argentinaTimeSingletonID = func() primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex("000000000000000000000001")
	if err != nil {
		return primitive.NilObjectID
	}
	return id
}()

/* UpsertArgentinaTimeDisplay guarda hora/día/mes/año en la colección time */
func UpsertArgentinaTimeDisplay(row models.ArgentinaTimeDisplay) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(timeCollection)

	doc := bson.M{
		"_id":       argentinaTimeSingletonID,
		"hora":      row.Hora,
		"fecha":     row.Fecha,
		"dia":       row.Dia,
		"mes":       row.Mes,
		"año":       row.Ano,
		"updatedat": time.Now(),
	}

	opts := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": argentinaTimeSingletonID}, doc, opts)
	if err != nil {
		return false, err
	}
	return true, nil
}
