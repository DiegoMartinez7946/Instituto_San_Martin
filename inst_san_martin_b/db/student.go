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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ensureStudentActiveField sync.Once

func ensureStudentActiveDefaults(ctx context.Context) {
	ensureStudentActiveField.Do(func() {
		db := config.MongoConnection.Database("san_martin")
		collection := db.Collection("student")
		_, err := collection.UpdateMany(ctx, bson.M{"active": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"active": true}})
		if err != nil {
			log.Println("ensure student active:", err)
		}
	})
}

/* GetStudentsDB get all students */
func GetStudentsDB() ([]*models.Student, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ensureStudentActiveDefaults(ctx)

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	var results []*models.Student

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cur, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Println(err.Error())
		return results, false
	}

	for cur.Next(ctx) {
		var row models.Student
		if err := cur.Decode(&row); err != nil {
			return results, false
		}
		results = append(results, &row)
	}
	_ = cur.Close(ctx)
	return results, true
}

/* InsertStudentDB insert one student */
func InsertStudentDB(s models.Student) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	row := bson.M{
		"name":          s.Name,
		"email":         s.Email,
		"phone":         s.Phone,
		"dni":           s.DNI,
		"address":       s.Address,
		"nivelaprobado": s.NivelAprobado,
		"degreeids":     s.DegreeIDs,
		"active":        s.Active,
		"createdat":     s.CreatedAt,
		"updatedat":     s.UpdatedAt,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar el alumno", err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil
}

/* UpdateStudentDB update student */
func UpdateStudentDB(s models.Student) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	filter := bson.M{"_id": s.ID}
	update := bson.M{
		"$set": bson.M{
			"name":          s.Name,
			"email":         s.Email,
			"phone":         s.Phone,
			"dni":           s.DNI,
			"address":       s.Address,
			"nivelaprobado": s.NivelAprobado,
			"degreeids":     s.DegreeIDs,
			"active":        s.Active,
			"updatedat":     s.UpdatedAt,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* FindStudentByDNIDB returns student if DNI exists (optional exclude id for update) */
func FindStudentByDNIDB(dni string, excludeID primitive.ObjectID) (models.Student, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	filter := bson.M{"dni": dni}
	if !excludeID.IsZero() {
		filter["_id"] = bson.M{"$ne": excludeID}
	}

	var st models.Student
	err := collection.FindOne(ctx, filter).Decode(&st)
	if err != nil {
		return st, false
	}
	return st, st.DNI != ""
}

/* UpdateStudentActiveDB solo estado activo/inactivo */
func UpdateStudentActiveDB(id primitive.ObjectID, active bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"active":    active,
			"updatedat": time.Now(),
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
