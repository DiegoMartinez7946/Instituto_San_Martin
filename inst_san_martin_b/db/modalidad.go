package db

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const modalidadCollection = "modalidad"

var modalidadSeedOnce sync.Once

/* EnsureModalidadSeeded inserta titular / provisional / suplente si la colección está vacía */
func EnsureModalidadSeeded() {
	modalidadSeedOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		collection := config.MongoConnection.Database("san_martin").Collection(modalidadCollection)
		n, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil || n > 0 {
			return
		}
		_, err = collection.InsertMany(ctx, []interface{}{
			bson.M{"codigo": "TITULAR"},
			bson.M{"codigo": "PROVISIONAL"},
			bson.M{"codigo": "SUPLENTE"},
		})
		if err != nil {
			log.Println("modalidad seed:", err)
		}
	})
}

/* GetModalidadesDB lista el catálogo */
func GetModalidadesDB() ([]*models.ModalidadDoc, bool) {
	EnsureModalidadSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(modalidadCollection)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil, false
	}
	defer cur.Close(ctx)
	var out []*models.ModalidadDoc
	for cur.Next(ctx) {
		var row models.ModalidadDoc
		if err := cur.Decode(&row); err != nil {
			return out, false
		}
		out = append(out, &row)
	}
	return out, true
}

/* GetModalidadByID returns document if exists */
func GetModalidadByID(id primitive.ObjectID) (models.ModalidadDoc, bool) {
	EnsureModalidadSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var row models.ModalidadDoc
	err := config.MongoConnection.Database("san_martin").Collection(modalidadCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&row)
	if err != nil || row.ID.IsZero() {
		return row, false
	}
	return row, true
}
