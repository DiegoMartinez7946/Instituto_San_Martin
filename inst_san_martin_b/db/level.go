package db

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
)

const levelCollection = "level"

var levelSeedOnce sync.Once

func defaultLevelRows() []interface{} {
	return []interface{}{
		bson.M{"nivel": "inicial", "jerarquia": 1},
		bson.M{"nivel": "primario", "jerarquia": 2},
		bson.M{"nivel": "secundario", "jerarquia": 3},
		bson.M{"nivel": "terciario", "jerarquia": 4},
		bson.M{"nivel": "universitario", "jerarquia": 4},
	}
}

/* EnsureLevelSeeded inserta filas por defecto nivel → orden si la colección está vacía */
func EnsureLevelSeeded() {
	levelSeedOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		db := config.MongoConnection.Database("san_martin")
		collection := db.Collection(levelCollection)

		n, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Println("level CountDocuments:", err)
			return
		}
		if n > 0 {
			return
		}
		_, err = collection.InsertMany(ctx, defaultLevelRows())
		if err != nil {
			log.Println("level InsertMany:", err)
		}
	})
}

/* GetLevelMap devuelve nivel (minúsculas) → valor jerarquia numérico */
func GetLevelMap() map[string]int {
	EnsureLevelSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(levelCollection)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("level Find:", err)
		return map[string]int{}
	}
	defer cur.Close(ctx)

	out := make(map[string]int)
	for cur.Next(ctx) {
		var row models.LevelHierarchy
		if err := cur.Decode(&row); err != nil {
			continue
		}
		if row.Nivel != "" {
			out[row.Nivel] = row.Jerarquia
		}
	}
	return out
}

/* GetLevelsDB devuelve todos los documentos de niveles ordenados por jerarquia y nivel */
func GetLevelsDB() ([]*models.LevelHierarchy, bool) {
	EnsureLevelSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(levelCollection)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("level Find:", err)
		return nil, false
	}
	defer cur.Close(ctx)

	var results []*models.LevelHierarchy
	for cur.Next(ctx) {
		var row models.LevelHierarchy
		if err := cur.Decode(&row); err != nil {
			return results, false
		}
		results = append(results, &row)
	}
	return results, true
}
