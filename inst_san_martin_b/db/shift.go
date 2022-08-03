// Package db provides ...
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
/* GetShiftsDB get the shifts from db */
func GetShiftsDB() ([]*models.Shift, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("shift")

	var results []*models.Shift

	condition := bson.M {  }
	optionsQuery := options.Find()
	optionsQuery.SetSort(bson.D {{ Key: "type", Value: -1}})

	shifts, err := collection.Find(ctx, condition, optionsQuery)
	if err != nil {
		log.Fatal(err.Error())
		return results, err
	}

	for shifts.Next(context.TODO()) {
		var row models.Shift
		err := shifts.Decode(&row)
		if err != nil {
			return results, err
		}
		results = append(results, &row)
	}

	return results, nil
}

/***************************************************************/
/***************************************************************/
/* InsertShiftDB insert one shift in db */
func InsertShiftDB(s models.Shift) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("shift")

	row := bson.M {
		"type": s.Type,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar el turno", err
	}
	
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil 
}

/***************************************************************/
/***************************************************************/
/* CheckExistShift check if shift already exists */
func CheckExistShift(typeShift string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("shift")

	typeShift = strings.ToUpper(typeShift)
	condition := bson.M {
		"type": typeShift,
	}

	var result models.Shift

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Type != "") {
		return result.ID.Hex(), true, nil
	}

	return "", false, err
}

/***************************************************************/
/***************************************************************/
/* UpdateShiftDB update the shift in the db */
func UpdateShiftDB(s models.Shift) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("shift")

	row := make(map[string]interface{})
	row["type"] = s.Type

	updateString := bson.M {
		"$set": row,
	}

	var idShift string
	idShift = s.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idShift)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/***************************************************************/
/***************************************************************/
/* DeleteShiftDB delete the academy shift from the db */
func DeleteShiftDB(IDShift string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("shift")

	objID, _ := primitive.ObjectIDFromHex(IDShift)

	condition := bson.M {
		"_id": objID,
	}

	_, err := collection.DeleteOne(ctx, condition)
	return err
}

/***************************************************************/
/***************************************************************/
/* GetShiftDB get the academy shift by id */
func GetShiftDB(IDShift string) (models.Shift, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("shift")

	objID, _ := primitive.ObjectIDFromHex(IDShift)

	condition := bson.M {
		"_id": objID,
	}

	var shift models.Shift

	err := collection.FindOne(ctx, condition).Decode(&shift)
	return shift, err
}
