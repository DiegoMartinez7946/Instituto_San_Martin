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

const jerarquiaCollection = "jerarquia"

var jerarquiaSeedOnce sync.Once

func defaultJerarquiaRows() []interface{} {
	return []interface{}{
		bson.M{"nivel": "inicial", "jerarquia": 1},
		bson.M{"nivel": "primario", "jerarquia": 2},
		bson.M{"nivel": "secundario", "jerarquia": 3},
		bson.M{"nivel": "terciario", "jerarquia": 4},
		bson.M{"nivel": "universitario", "jerarquia": 4},
	}
}

/* EnsureJerarquiaSeeded inserts default nivel → jerarquia rows if the collection is empty */
func EnsureJerarquiaSeeded() {
	jerarquiaSeedOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		db := config.MongoConnection.Database("san_martin")
		collection := db.Collection(jerarquiaCollection)

		n, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Println("jerarquia CountDocuments:", err)
			return
		}
		if n > 0 {
			return
		}
		_, err = collection.InsertMany(ctx, defaultJerarquiaRows())
		if err != nil {
			log.Println("jerarquia InsertMany:", err)
		}
	})
}

/* GetJerarquiaMap returns nivel (lowercase) → jerarquia number */
func GetJerarquiaMap() map[string]int {
	EnsureJerarquiaSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(jerarquiaCollection)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("jerarquia Find:", err)
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

/* GetJerarquiasDB returns all hierarchy documents sorted by jerarquia then nivel */
func GetJerarquiasDB() ([]*models.LevelHierarchy, bool) {
	EnsureJerarquiaSeeded()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(jerarquiaCollection)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("jerarquia Find:", err)
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
