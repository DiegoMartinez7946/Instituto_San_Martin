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

const tituloHabilitanteCollection = "titulo_habilitante"

var tituloHabilitanteSeedOnce sync.Once

/* EnsureTituloHabilitanteSeeded inserts SI/NO si la colección está vacía */
func EnsureTituloHabilitanteSeeded() {
	tituloHabilitanteSeedOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		collection := config.MongoConnection.Database("san_martin").Collection(tituloHabilitanteCollection)
		n, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil || n > 0 {
			return
		}
		_, err = collection.InsertMany(ctx, []interface{}{
			bson.M{"codigo": "SI"},
			bson.M{"codigo": "NO"},
		})
		if err != nil {
			log.Println("titulo_habilitante seed:", err)
		}
	})
}

/* GetTitulosHabilitantesDB lista el catálogo */
func GetTitulosHabilitantesDB() ([]*models.TituloHabilitante, bool) {
	EnsureTituloHabilitanteSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(tituloHabilitanteCollection)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil, false
	}
	defer cur.Close(ctx)
	var out []*models.TituloHabilitante
	for cur.Next(ctx) {
		var row models.TituloHabilitante
		if err := cur.Decode(&row); err != nil {
			return out, false
		}
		out = append(out, &row)
	}
	return out, true
}

/* GetTituloHabilitanteByID returns document if exists */
func GetTituloHabilitanteByID(id primitive.ObjectID) (models.TituloHabilitante, bool) {
	EnsureTituloHabilitanteSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var row models.TituloHabilitante
	err := config.MongoConnection.Database("san_martin").Collection(tituloHabilitanteCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&row)
	if err != nil || row.ID.IsZero() {
		return row, false
	}
	return row, true
}
