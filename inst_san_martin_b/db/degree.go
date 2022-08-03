package db

import (
	"context"
	"log"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/***************************************************************/
/***************************************************************/
/* GetDegreesDB get the degrees from db */
func GetDegreesDB() ([]*models.Degree, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	var results []*models.Degree

	condition := bson.M {  }
	optionsQuery := options.Find()
	optionsQuery.SetSort(bson.D {{ Key: "name", Value: 1}, { Key: "active", Value: 1}})

	degrees, err := collection.Find(ctx, condition, optionsQuery)
	if err != nil {
		log.Fatal(err.Error())
		return results, false
	}

	for degrees.Next(context.TODO()) {
		var row models.Degree
		err := degrees.Decode(&row)
		if err != nil {
			return results, false
		}
		results = append(results, &row)
	}

	return results, true
}

/***************************************************************/
/***************************************************************/
/* InsertDegreeDB insert one degree in db */
func InsertDegreeDB(d models.Degree) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	row := bson.M {
		"name": d.Name,
		"active": d.Active,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar la carrera", err
	}
	
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil 
}

/***************************************************************/
/***************************************************************/
/* CheckExistDegree check if degree already exists */
func CheckExistDegree(nameDegree string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	condition := bson.M {
		"name": nameDegree,
	}

	var result models.Degree

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Name != "") {
		return result.ID.Hex(), true, nil
	}

	return "", false, err
}

/***************************************************************/
/***************************************************************/
/* UpdateDegreeDB update the degree in the db */
func UpdateDegreeDB(d models.Degree) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	row := make(map[string]interface{})
	row["name"] = d.Name

	updateString := bson.M {
		"$set": row,
	}

	var idDegree string
	idDegree = d.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idDegree)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateStatusDegreeDB update the degree in the db */
func UpdateStatusDegreeDB(d models.Degree) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	row := make(map[string]interface{})
	row["active"] = d.Active

	updateString := bson.M {
		"$set": row,
	}

	var idDegree string
	idDegree = d.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idDegree)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}
