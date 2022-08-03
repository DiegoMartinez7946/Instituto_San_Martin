package db

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/***************************************************************/
/***************************************************************/
/* GetTestTypesDB get the test types from db */
func GetTestTypesDB() ([]*models.TestType, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("test_type")

	var results []*models.TestType

	condition := bson.M {  }
	optionsQuery := options.Find()
	optionsQuery.SetSort(bson.D {{ Key: "type", Value: -1}})

	testTypes, err := collection.Find(ctx, condition, optionsQuery)
	if err != nil {
		log.Fatal(err.Error())
		return results, false
	}

	for testTypes.Next(context.TODO()) {
		var row models.TestType
		err := testTypes.Decode(&row)
		if err != nil {
			return results, false
		}
		results = append(results, &row)
	}

	return results, true
}

/***************************************************************/
/***************************************************************/
/* InsertTestTypeDB insert one test type in db */
func InsertTestTypeDB(t models.TestType) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("test_type")

	row := bson.M {
		"type": t.Type,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar el tipo de examen", err
	}
	
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil 
}

/***************************************************************/
/***************************************************************/
/* CheckExistTestType check if test type already exists */
func CheckExistTestType(typeTestType string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("test_type")

	typeTestType = strings.ToUpper(typeTestType)
	condition := bson.M {
		"type": typeTestType,
	}

	var result models.TestType

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Type != "") {
		return result.ID.Hex(), true, nil
	}

	return "", false, err
}

/***************************************************************/
/***************************************************************/
/* UpdateTestTypeDB update the test type in the db */
func UpdateTestTypeDB(t models.TestType) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("test_type")

	row := make(map[string]interface{})
	row["type"] = t.Type

	updateString := bson.M {
		"$set": row,
	}

	var idTestType string
	idTestType = t.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idTestType)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/***************************************************************/
/***************************************************************/
/* DeleteTestTypeDB delete the test type from the db */
func DeleteTestTypeDB(IDTestType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("test_type")

	objID, _ := primitive.ObjectIDFromHex(IDTestType)

	condition := bson.M {
		"_id": objID,
	}

	_, err := collection.DeleteOne(ctx, condition)
	return err
}

/***************************************************************/
/***************************************************************/
/* GetTestTypeDB get the test type by id */
func GetTestTypeDB(IDTestType string) (models.TestType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("test_type")

	objID, _ := primitive.ObjectIDFromHex(IDTestType)

	condition := bson.M {
		"_id": objID,
	}

	var testType models.TestType

	err := collection.FindOne(ctx, condition).Decode(&testType)
	return testType, err
}
